package main
//SOME ADDED TEST CASES . Also please note that the semver library being used is the one used in main.go as well for
//obvious compatibility
import (
	"testing"
	"fmt"
	"github.com/Masterminds/semver"
)
func myNewVersion (version string) *semver.Version {
  result,err := semver.NewVersion(version)
  if(err!=nil){
    fmt.Errorf("Error parsing segment %s",err)
    return nil
  }
  return result
}
func stringToVersionSlice(stringSlice []string) []*semver.Version {
	versionSlice := make([]*semver.Version, len(stringSlice))
	for i, versionString := range stringSlice {
		versionSlice[i] = myNewVersion(versionString)
	}
	return versionSlice
}

func versionToStringSlice(versionSlice []*semver.Version) []string {
	stringSlice := make([]string, len(versionSlice))
	for i, version := range versionSlice {
		stringSlice[i] = version.String()
	}
	return stringSlice
}

func TestLatestVersions(t *testing.T) {
	testCases := []struct {
		versionSlice   []string
		expectedResult []string
		minVersion     *semver.Version
	}{
		{
			versionSlice:   []string{"1.8.11", "1.9.6", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6", "1.8.11"},
			minVersion:     myNewVersion("1.8.0"),
		},
		{
			versionSlice:   []string{"1.8.11", "1.9.6", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6"},
			minVersion:     myNewVersion("1.8.12"),
		},
		{
			versionSlice:   []string{"1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1"},
			minVersion:     myNewVersion("1.10.0"),
		},
		{
			versionSlice:   []string{"2.2.1", "2.2.0"},
			expectedResult: []string{"2.2.1"},
			minVersion:     myNewVersion("2.2.1"),
		},
		// Implement more relevant test cases here, if you can think of any
    		//All versions lower than min Version so return empty slice
    {
			versionSlice:   []string{"1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{},
			minVersion:     myNewVersion("1.11.0"),
		},
    		// version slice empty
    {
			versionSlice:   []string{},
			expectedResult: []string{},
			minVersion:     myNewVersion("1.2.0"),
		},
    		//Prerelease vs No pre release
    {
			versionSlice:   []string{"1.8.11", "1.9.6", "1.9.5", "1.8.10","1.11.0-beta","1.11.0", "1.10.1-beta", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.11.0","1.10.1-beta", "1.9.6",},
			minVersion:     myNewVersion("1.9.0"),
		},
    		//Major version number difference,not only minor number difference
    {
			versionSlice:   []string{"1.8.11", "1.9.6", "1.9.5", "1.8.10","1.11.0-beta","2.0.0","2.0.1","2.1.0","2.1.0-beta","1.11.0", "1.10.1-beta", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"2.1.0","2.0.1","1.11.0","1.10.1-beta", "1.9.6",},
			minVersion:     myNewVersion("1.9.0"),
		}, 
	}

	test := func(versionData []string, expectedResult []string, minVersion *semver.Version) {
		stringSlice := versionToStringSlice(LatestVersions(stringToVersionSlice(versionData), minVersion))
		if (len(stringSlice)!=len(expectedResult)){
			t.Errorf("Received %s ,expected %s",stringSlice,expectedResult)
			return
		}
		for i, versionString := range stringSlice {
			if versionString != expectedResult[i] {
				t.Errorf("Received %s, expected %s", stringSlice, expectedResult)
				return
			}
		}
	}

	for _, testValues := range testCases {
		test(testValues.versionSlice, testValues.expectedResult, testValues.minVersion)
	}
}
