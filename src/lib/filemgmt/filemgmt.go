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
// TAG DETAILS:
// 0 -> File will be split for testing 1 node
// n > 0 -> File will be split for testing on n nodes
// -------------------
// INSTRUCTIONS:
// -------------------
// Call filemgmt.SplitFile("{your filename}", {tag})
func SplitFile(filename string, tag int) {
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
func ProcessLocalFiles(localFileSys *utility.FileSys) {
	locFiles, err := ioutil.ReadFile(consts.DirPath + "/localFiles.json")
	utility.CheckError(err)
	var filePaths utility.FilePath
	files := make([]utility.File, 0)
	filePaths.Files = files
	err = json.Unmarshal(locFiles, &filePaths)
	utility.CheckError(err)
	fmt.Println("======================    PROCESSING LOCAL FILES INTO FILE SYSTEM    =====================")
	fmt.Println("==========================================================================================")
	for index, value := range filePaths.Files {
		colorprint.Info("---------------------------------------------------------------------------")
		fmt.Println((index + 1), ">> PROCESSING:", value.Name, "at "+value.Path)
		substrind := strings.Index(value.Name, ".")
		substr := value.Name[:substrind]
		var vidmap map[int]utility.VidSegment
		vidmap = make(map[int]utility.VidSegment)
		// var vidBytes []byte
		for i := 0; i < len(value.SegsAvail); i++ {
			pathname := value.Path + substr + "_" + strconv.Itoa(int(value.SegsAvail[i]))
			fmt.Printf("\rProcessing segment %s for %s out of %d segments", strconv.Itoa(int(value.SegsAvail[i])), value.Name, value.SegNums)
			dat, err := ioutil.ReadFile(pathname)
			utility.CheckError(err)
			var vidSeg utility.VidSegment
			err = json.Unmarshal(dat, &vidSeg)
			utility.CheckError(err)
			vidmap[int(value.SegsAvail[i])] = vidSeg
			// for j := 0; j < len(vidSeg.Body); j++ {
			// 	vidBytes = append(vidBytes, vidSeg.Body[j])
			// }
		}
		vid := utility.Video{
			Name:      value.Name,
			SegNums:   value.SegNums,
			SegsAvail: value.SegsAvail,
			Segments:  vidmap,
		}
		localFileSys.Lock()
		colorprint.Alert("\nLOCKING FILESYSTEM")
		localFileSys.Files[value.Name] = vid
		localFileSys.Unlock()
		colorprint.Alert("\nUNLOCKING FILESYSTEM")
		// err := ioutil.WriteFile(consts.DirPath+"/saved/"+value.Name, vidBytes, 0777)
		// utility.CheckError(err)
		colorprint.Debug("Locally available segments for " + value.Name + " saved into the filesystem")
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

// printFileSysContents(localFileSys *utility.FileSys)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Print the filesystem contents
func PrintFileSysContents(localFileSys *utility.FileSys) {
	colorprint.Debug("READING CONTENTS OF FILE SYSTEM NO. " + strconv.Itoa(localFileSys.Id))
	localFileSys.Lock()
	colorprint.Warning("======================================================================================================")
	for index, value := range localFileSys.Files {

		colorprint.Warning("--- FILE: " + index + " ---")
		var altSegsAvail []string
		if len(value.SegsAvail) > 5 {
			for i := 0; i < len(value.SegsAvail); i++ {
				if i <= 5 {
					altSegsAvail = append(altSegsAvail, strconv.Itoa(int(value.SegsAvail[i])))
				} else {
					altSegsAvail = append(altSegsAvail, "...")
					altSegsAvail = append(altSegsAvail, strconv.Itoa(int(value.SegsAvail[len(value.SegsAvail)-1])))
					break
				}
			}
		}
		fmt.Println("FILENAME:", value.Name, "\nTOTAL SEGMENTS:", value.SegNums, "\nSEGMENTS AVAILABLE:", altSegsAvail)
	}
	localFileSys.Unlock()
	colorprint.Warning("======================================================================================================")
}

// // addVidSegIntoFileSys(filename string, vidSeg utility.VidSegment) {
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Adds the video segment info into localFiles.json. Call this function after adding a new video segment to the file system
func AddVidSegIntoFileSys(filename string, segNums int64, vidSeg utility.VidSegment, localFileSys *utility.FileSys) {
	localFileSys.RLock()
	colorprint.Debug("LOCK")
	_, ok := localFileSys.Files[filename]
	localFileSys.RUnlock()
	if ok {
		localFileSys.Files[filename].Segments[vidSeg.Id] = vidSeg
	} else {
		var vidmap map[int]utility.VidSegment
		var arr []int64
		arr = append(arr, int64(vidSeg.Id))
		vidmap[vidSeg.Id] = vidSeg
		vid := utility.Video{
			Name:      filename,
			SegNums:   segNums,
			SegsAvail: arr,
			Segments:  vidmap,
		}
		localFileSys.Lock()
		localFileSys.Files[filename] = vid
		localFileSys.Unlock()

	}
	colorprint.Debug("UNLOCK")
}

// // addVidSegIntoFileSysJSON(filename string, vidSeg utility.VidSegment) {
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Adds the video segment info into localFiles.json. Call this function after adding a new video segment to the file system
func addVidSegIntoFileSysJSON(filename string, path string, segNums int, vidSeg utility.VidSegment) {
	vidSegId := int64(vidSeg.Id)
	locFiles, err := ioutil.ReadFile(consts.DirPath + "/localFiles.json")
	utility.CheckError(err)
	var filePaths utility.FilePath
	files := make([]utility.File, 0)
	filePaths.Files = files
	err = json.Unmarshal(locFiles, &filePaths)
	utility.CheckError(err)
	var filefound bool = false
	for _, value := range filePaths.Files {
		if value.Name == filename {
			value.SegsAvail = append(value.SegsAvail, vidSegId)
			filefound = true
			break
		}
	}
	if !filefound {
		var file utility.File
		file.Name = filename
		file.Path = path
		file.SegNums = int64(segNums)
		file.SegsAvail = append(file.SegsAvail, vidSegId)
	}

	dat, err := json.Marshal(filePaths)
	utility.CheckError(err)
	utility.SaveFileInfoToJson(dat)
}
