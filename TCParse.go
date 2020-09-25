package main

import "fmt"
import "bufio"
import "os"
import "path"
import "log"
import "strings"

//import "strings"

func main() {
	var errLines []string
	var eDate string
	var eTime string
	var eLvl string
	var eMsg string
	var oName string
	var next2 int
	var isOErr bool
	var oErr string
	var oFound bool
	var oTestCnt int
	var lineCount int
	var verbose string
	var haveFile bool

	type oMessages struct {
		oErr    string
		objName string
		eCount  int
	}

	if len(os.Args) < 5 {
		fmt.Println("Usage:  TomcatLogParse RemoteTomcatLog LocalWorkingDir RemoteUser RemoteServer verbosity[s,v,vv]")
		fmt.Println("Example: TOmcatLogParse /usr/local/tomcat/logs/catalina.log c:/Users/work/ user server s")
		fmt.Println("Verbosity Level: s = Summary, v = Verbose, vv = Very Verbose")
		os.Exit(1)
	}
	verbose = os.Args[5]
	catFname := os.Args[1]
	if strings.Compare(verbose, "v") == 0 || strings.Compare(verbose, "vv") == 0 {
		fmt.Printf("Catalina Filename: %s\n", catFname)
	}

	laFname := os.Args[2] + path.Base(catFname)
	if strings.Compare(verbose, "v") == 0 || strings.Compare(verbose, "vv") == 0 {
		fmt.Printf("Working Directory Filename: %s\n", laFname)
	}

	haveFile = GetDataFile(catFname, laFname, os.Args[4], os.Args[3])

	if !haveFile {
		log.Fatal("Could Not Retrieve Log")
	} else {
		fmt.Printf("Received File: %s\nFrom: %s\n", catFname, os.Args[4])
	}

	f, err := os.Open(laFname)
	if err != nil {
		log.Fatalf("Error opening catalina file: %s\n", laFname)
	}
	defer f.Close()

	cinFile := bufio.NewScanner(f)

	oMCnt := []oMessages{}
	next2 = 0
	lineCount = 0
	fmt.Println("Oracle Errors")
	fmt.Println("-------------")
	for cinFile.Scan() {
		lineCount += 1
		errLines = append(errLines, cinFile.Text())
		if next2 <= 0 {
			eDate, eTime, eLvl, eMsg = parceTCLine(cinFile.Text())
			if eDate != "NO" {
				if strings.Compare(verbose, "v") == 0 || strings.Compare(verbose, "vv") == 0 {
					fmt.Printf("Date: %s Time: %s\n", eDate, eTime)
					fmt.Printf("Level: %s Message: %s\n", eLvl, eMsg)
				}
				isOErr, oErr = isOracleError(eMsg)
				if strings.Compare(verbose, "v") == 0 || strings.Compare(verbose, "vv") == 0 {
					if isOErr {
						fmt.Printf("Oracle Error: %s\n", oErr)
					}
				}
				next2 = 2
			}
		} else {
			if next2 == 2 {
				oName = getOname(cinFile.Text())
				if strings.Compare(verbose, "v") == 0 || strings.Compare(verbose, "vv") == 0 {
					fmt.Printf("Object Name: %s\n", oName)
				}
			}
			if strings.Compare(verbose, "v") == 0 || strings.Compare(verbose, "vv") == 0 {
				fmt.Printf("Error[%d]: %s\n", next2, cinFile.Text())
			}
			next2--
		}
		if len(oName) > 0 && len(oErr) > 0 {
			oFound = false
			for i := 0; i < len(oMCnt); i++ {
				if oMCnt[i].oErr == oErr && oMCnt[i].objName == oName {
					oMCnt[i].eCount++
					oFound = true
				}
			}
			if !oFound { //No Obj Found
				oMCnt = append(oMCnt, oMessages{oErr: oErr, objName: oName, eCount: 1})
			}
			oName = ""
			oErr = ""
		}

	}
	if err := cinFile.Err(); err != nil {
		fmt.Println("Error reading catalina file")
	}

	for i := 0; i < len(oMCnt); i++ {
		fmt.Printf("Object: %s Error: %s Count: %d\n", oMCnt[i].objName, oMCnt[i].oErr, oMCnt[i].eCount)
	}
	oTestCnt += 1
	fmt.Printf("File Processed with %d line(s).\n", lineCount)
	fmt.Println("TomcatLogParse Finished")

}
