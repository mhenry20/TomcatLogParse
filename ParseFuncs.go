package main

import "strings"

func parceTCLine(line string) (eDate string, eTime string, eLevel string, eMsg1 string) {
	var tokArray []string
	//fmt.Printf("Line: %s\n", line)
	tokArray = strings.Fields(line)
	if len(tokArray) >= 7 {
		//fmt.Printf("Tok[3] = %s\n",tokArray[3])
		if strings.HasPrefix(tokArray[3], "[ajp-nio") { //First Message Line
			if strings.HasPrefix(tokArray[6], "ORA-") {
				return tokArray[0], tokArray[1], tokArray[2], tokArray[6]
			}
		}
	}

	return "NO", "NO", "NO", "NO"
}

func isOracleError(eMsg string) (isOra bool, eNumb string) {
	if strings.HasPrefix(eMsg, "ORA-") {
		isOra = true
		eNumb = strings.Split(eMsg, "-")[1]
		eNumb = strings.Replace(eNumb, ":", "", -1)
		return isOra, eNumb
	} else {
		return false, ""
	}

}

func getOname(eMsg string) (oname string) {
	var firstInd, lastInd int
	if strings.HasPrefix(eMsg, "PLS-") {
		firstInd = strings.Index(eMsg, "'") + 1
		lastInd = strings.LastIndex(eMsg, "'")
		return eMsg[firstInd:lastInd]
	}
	return ""
}
