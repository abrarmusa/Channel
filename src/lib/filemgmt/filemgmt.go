package filemgmt

import (
	"../../consts"
	"../colorprint"
	"../utility"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// SplitFile(filename string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Converts a file into several json encoded segments
func SplitFile(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	utility.CheckError(err)
	var eightBSeg []byte
	var vidSeg utility.VidSegment
	counter := 1
	foldername := procName(filename)
	colorprint.Alert("VIDEO HAS " + strconv.Itoa(len(bytes)) + " bytes. These will be divided up into " + strconv.Itoa(len(bytes)/consts.Bytecount) + " segments.")
	for index, element := range bytes {
		// colorprint.Debug("INDEX: " + strconv.Itoa(index) + " COUNTER:" + strconv.Itoa(counter))
		eightBSeg = append(eightBSeg, element)
		ident := ((index / consts.Bytecount) + 1)
		// colorprint.Debug(strconv.Itoa(index) + " " + strconv.Itoa(ident))
		if counter == consts.Bytecount {
			counter = 0
			vidSeg = utility.VidSegment{
				Id:   ident,
				Body: eightBSeg,
			}
			data, err := json.Marshal(vidSeg)
			utility.CheckError(err)
			eightBSeg = []byte{}

			if _, err := os.Stat(consts.DirPath + consts.LocalPath + foldername); os.IsNotExist(err) {
				colorprint.Warning("Creating folder " + consts.DirPath + consts.LocalPath + foldername)
				err := os.MkdirAll(consts.DirPath+consts.LocalPath+foldername, 0777)
				utility.CheckError(err)
			}
			writeToFileHelper(foldername, ident, data)
		}
		counter++
	}

}

// func ProcessLocalFiles(fileSys *utility.FileSys)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Looks into the json for a filename and processes the segments into the filesystem.
// NOTE: A POINTER TO THE LOCAL FILESYSTEM MUST BE INPUT
func ProcessLocalFiles() {
	locFiles, err := ioutil.ReadFile(consts.DirPath + "/localFiles.json")
	fmt.Println("RED")
	utility.CheckError(err)
	var filePaths utility.FilePath
	files := make([]utility.File, 0)
	filePaths.Files = files
	err = json.Unmarshal(locFiles, &filePaths)
	utility.CheckError(err)
	fmt.Println("========================    PROCESSING LOCAL FILES FOR SHARING    ========================")
	fmt.Println("==========================================================================================")
	for index, value := range filePaths.Files {
		colorprint.Info("---------------------------------------------------------------------------")
		fmt.Println((index + 1), ">> PROCESSING:", value.Name, "at "+value.Path)
		substrind := strings.Index(value.Name, ".")
		substr := value.Name[:substrind]
		var vidBytes []byte
		for i := 0; i < len(value.SegsAvail); i++ {
			pathname := value.Path + substr + "_" + strconv.Itoa(value.SegsAvail[i])
			fmt.Printf("\rProcessing segment %s for %s", strconv.Itoa(value.SegsAvail[i]), value.Name)
			dat, err := ioutil.ReadFile(pathname)
			utility.CheckError(err)
			var vidSeg utility.VidSegment
			err = json.Unmarshal(dat, &vidSeg)
			utility.CheckError(err)
			for j := 0; j < len(vidSeg.Body); j++ {
				vidBytes = append(vidBytes, vidSeg.Body[j])
			}
			utility.CheckError(err)
		}
		fmt.Println()
		err := ioutil.WriteFile(consts.DirPath+"/saved/"+value.Name, vidBytes, 0777)
		utility.CheckError(err)
		colorprint.Info("---------------------------------------------------------------------------")

	}
	fmt.Println("===============================    PROCESSING COMPLETE    ================================\n\n\n")
}

// writeToFileHelper()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// This method writes the downloaded file into a file of its own in the filesys/downloaded folder
func writeToFileHelper(foldername string, ident int, data []byte) {
	str := consts.DirPath + consts.LocalPath + foldername + "/" + foldername + "_" + strconv.Itoa(ident)
	err := ioutil.WriteFile(str, data, 0777)
	utility.CheckError(err)
}

// procName(filename string) string
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Processes the filename into the appropriate folder name for the segments to be stored
func procName(filename string) string {
	str := strings.Split(filename, "/")
	foldername := str[len(str)-1]
	dotindex := strings.Index(foldername, ".")
	return filename[:dotindex]
}

// saveVidSeg(filename string, vidSeg utility.VidSegment) string
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Processes the filename into the appropriate folder name for the segments to be stored
func saveVidSegIntoFileSys(filename string, vidSeg utility.VidSegment) {
	// TODO
	// TODO
	// TODO
	// TODO
}

// addVidSegIntoFileSysJSON(filename string, vidSeg utility.VidSegment) {
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Adds the video segment info into localFiles.json
func addVidSegIntoFileSysJSON(filename string, ident int) {
	// TODO
	// TODO
	// TODO
	// TODO
}

// addVidSegIntoFileSysJSON(filename string, vidSeg utility.VidSegment) {
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Adds the video segment info into localFiles.json
func addVidSegIntoFileSysJSON(filename string, ident int) {
	// TODO
	// TODO
	// TODO
	// TODO
}

// // processLocalVideosIntoFileSys()
// // --------------------------------------------------------------------------------------------
// // DESCRIPTION:
// // -------------------
// // This method loads up a local json file to see which files are available in the local file system. Once
// // the read has been completed, the files are then processed into the utility.utility.FileSys map accordingly
// func processLocalVideosIntoFileSys() {
// 	locFiles, err := ioutil.ReadFile(consts.DirPath + "/localFiles.json")
// 	utility.CheckError(err)
// 	files := make([]utility.File, 0)

// 	filePaths.Files = files
// 	err = json.Unmarshal(locFiles, &filePaths)
// 	utility.CheckError(err)
// 	// Initialize local file system
// 	localFileSys = utility.FileSys{
// 		Id:    1,
// 		Files: make(map[string]utility.Video),
// 	}
// 	fmt.Println("========================    PROCESSING LOCAL FILES FOR SHARING    ========================")
// 	fmt.Println("==========================================================================================")
// 	for index, value := range filePaths.Files {

// 		dat, err := ioutil.ReadFile(value.Path)
// 		utility.CheckError(err)
// 		colorprint.Info("---------------------------------------------------------------------------")
// 		colorprint.Info(strconv.Itoa(index+1) + ": Processing " + value.Name + " at " + value.Path + " with " + strconv.Itoa(len(dat)/consts.Bytecount) + " segments.")
// 		segsAvail, vidMap := convByteArrayToSeg(dat)

// 		vid := utility.Video{
// 			Name:      value.Name,
// 			SegNums:   int64(len(dat) / consts.Bytecount),
// 			SegsAvail: segsAvail,
// 			Segments:  vidMap,
// 		}
// 		localFileSys.Lock()
// 		localFileSys.Files[value.Name] = vid
// 		localFileSys.Unlock()
// 		colorprint.Info("Completed Processing " + value.Name + " at " + value.Path)
// 		colorprint.Info("---------------------------------------------------------------------------")

// 	}
// 	fmt.Println("===============================    PROCESSING COMPLETE    ================================\n\n\n")

// }

// // convByteArrayToSeg(bytes []byte) ([]int64, map[int]utility.VidSegment)
// // --------------------------------------------------------------------------------------------
// // DESCRIPTION:
// // -------------------
// // Converts the byte array from a video files into utility.Video Segments.
// func convByteArrayToSeg(bytes []byte) ([]int64, map[int]utility.VidSegment) {
// 	vidmap := make(map[int]utility.VidSegment)
// 	var segsAvail []int64
// 	var vidSeg utility.VidSegment
// 	var eightBSeg []byte
// 	counter, counter2, counter3 := 1, 1, 1
// 	progstr := "="
// 	blen := len(bytes)
// 	altc := (blen / int(consts.Factor))
// 	for index, element := range bytes {
// 		eightBSeg = append(eightBSeg, element)
// 		if counter == consts.Bytecount {
// 			counter = 0
// 			vidSeg = utility.VidSegment{
// 				Id:   ((index % consts.Bytecount) + 1),
// 				Body: eightBSeg,
// 			}
// 			vidmap[((index / consts.Bytecount) + 1)] = vidSeg
// 			segsAvail = append(segsAvail, int64(((index / consts.Bytecount) + 1)))
// 			eightBSeg = []byte{}
// 		}
// 		counter++
// 		counter2++
// 		counter3++
// 		if counter2 == altc {
// 			progstr += "~"
// 			fmt.Printf("\r|%s|  - %d%%", progstr, ((counter3*100)/blen + 1))
// 			counter2 = 0
// 		}
// 	}
// 	fmt.Println()
// 	return segsAvail, vidmap
// }
