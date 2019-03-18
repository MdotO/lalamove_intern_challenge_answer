package main
// Used a Masterminds/semver instead of the given one since this one gave auseful utility function sort.
import (
	"context"
	"fmt"
	"os"
	"log"
	"bufio"
	"sort"
	"strings"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
)


// A struct to hold related info unit from each line read from file
type RepoInfo struct{
	name string
	repo string
	minVersion string
} 


//Function that distinguishes between two version ,whether the first one has a greater major or minor version or both
func GreaterMajorOrMinor(a *semver.Version,b *semver.Version) bool{
	if (a.Major() > b.Major()||
	    (a.Major() == b.Major()&&
	     a.Minor() > b.Minor())){
		return true
	}
	return false
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor //versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	//Sort the released versions(O(nlogn)) and then perform checks on them in linear time   
	sort.Sort(semver.Collection(releases))
	x := len(releases)
	//base condition
	if(x==0){
		return versionSlice
	}
 	//Being sorted ,just compare and append when cosnecutive versions do differ by either major/minor versions 
	for i,r := range releases{
		if (i!=0){
			//r is current version in consideration i.e releases[i]
			if (GreaterMajorOrMinor(r,releases[i-1]) &&
			    releases[i-1].Compare(minVersion)>=0){
				versionSlice = append(versionSlice,releases[i-1])
			}
		}
	}
	//Append the last element always!!!
	if (releases[x-1].Compare(minVersion)>=0){
		versionSlice = append(versionSlice,releases[x-1])
	}
  	//reversing the obtained list to get the desired format
   	for i, j := 0, len(versionSlice)-1; i < j; i, j = i+1, j-1 {
        	versionSlice[i], versionSlice[j] = versionSlice[j], versionSlice[i]
    	}
		return versionSlice
	}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	//NOTE: Importantly, the file format {repository,minVersion as line1 with repository owner,name & minversion in ALL NEXT 
	//LINES TILL EOF }
	file,err := os.Open(os.Args[1])
  	if(err!=nil){
    	log.Fatal(err)
  	}
  	defer file.Close()
	//Start scanning
  	scanner := bufio.NewScanner(file)

  	//First line just column headers,Nothing else so just read and discard
  	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
  	var Repos[] RepoInfo 
  	scanner.Scan()
	
	//Sort out each line of info
  	for scanner.Scan(){
    		text :=scanner.Text()
    		var temp RepoInfo 
    		temp.name = (strings.Split((strings.Split(text,","))[0],"/"))[0]
    		temp.repo = (strings.Split((strings.Split(text,","))[0],"/"))[1]
    		temp.minVersion = strings.Split(text,",")[1]
		//Append to the list
    		Repos = append(Repos,temp)
  	}
  	length_repos := len(Repos)
	//Use loop for each API call since NO SUCH GITHUB API EXISTS that can take a list of owners+repos and give their releases
	for x :=0; x<length_repos;x++{
		releases, _, err := client.Repositories.ListReleases(ctx, Repos[x].name, Repos[x].repo, opt)
		if err != nil {
			panic(err) // is this really a good way?
		}
		minVersion,_ := semver.NewVersion(Repos[x].minVersion)
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i],_ = semver.NewVersion(versionString)
		}

		versionSlice := LatestVersions(allReleases, minVersion)
		//Print in format(fmt.Printf()) as described in specification. 
		fmt.Printf("latest versions of %s/%s: %s", Repos[x].name,Repos[x].repo,versionSlice)
		//New line print if and only if not reached the end 
  		if(x<length_repos-1){
    			fmt.Println("")
		}
	}
}
