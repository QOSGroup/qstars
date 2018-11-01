# SDK development for iOS and Android
How to use Gomobile for library development for iOS and Android

## Environment prerequisites:
* OSX
* xcode
* golang 1.5 +
* Android SDK
* Andorid NDK
* JDK

## Gomobile installation 
Let’s start with gomobile installation:
```
go get -v golang.org/x/mobile/cmd/gomobile
```
Note: in OS X you need to have installed Xcode Command Line Tools. Then you need to initialize gomobile, this can be done one time in any work directory.
```
gomobile init
```
Note: this command might take several minutes.

## Android SDK environment setting
Download and unpack Android SDK to home directory, for example, `~/android-sdk`, and make the following command for API installation.
```
~/android-sdk/tools/android sdk
```
Then you need to set environment variable:
```
export ANDROID_HOME=$HOME"/android-sdk"
```
So far the environment for the library development and building is ready.

## Shared GO-code for Android and iOS
The same code can be used for the future compilation for Android and iOS. Building of such cross platform code has its own constraints. As for now we can use only certain set of data types. We need to take it into consideration when developing application in Go. Let’s review in more detail the supported types:
* int and float;
* string and boolean;
* []byte
* function has to return only supported types, it may not return the result, it may return one or two types wherein the second type should be an error;
* interfaces could be used if they are exported to files of any supported type;
* struct type, only in case all fields meet the constraints.

So, if the type is not supported by gomobile bind command, you’ll see the similar error:
```
panic: unsupported basic seqType: uint64
```
It’s obvious that the set of supported types is very limited, but this is enough for the SDK implementation.

## Fetch the source Go code from Github repository
The source code is based on Golang and the repository for the project is <https://github.com/QOSGroup/qstars>. 
It could be git cloned and downloaded to your local repo at first via:
 ```
 git clone https://github.com/QOSGroup/qstars.git
 ```
The source code is under the folder `qstars/stub/ios/` with name of `starsdk.go`

Note: It is importtant to fetch the repository and corresponding packages this project, i.e. `qstars` depends. Otherwise, the gomobile would encounter buid failure.


## Building and import to Java/Objective-C/Swift
Gobind generates target language (Java, Objective-C or Swift) bindings for each exported symbol in a Go package.
The code is generated automatically and is packed by `gomobile bind` command. More details you can find here <https://golang.org/x/mobile/cmd/gomobile>.

Let’s start with flag -target that defines platform for generation. Here is an example for Android:

```
gomobile bind --target=android .
```
This command will generate `.aar` file from the current code. To import that file to Android Studio is pretty simple:
* File ➤New ➤New Module ➤Import .JAR or .AAR package
* File ➤Project Structure ➤app ➤Dependencies ➤Add Module Dependency
* Add import: import go.logpackermobilesdk.Logpackermobilesdk

Note: In Java the name of the package for import always starts with go.

Similar command is used for building Objective-C/Swift code.
```
gomobile bind --target=ios .
```
The folder `.framework` will be created in the current repository.

This works for both Objective-C and Swift. Transfer `.framework` folder to Xcode’s file browser and add import to project:
```
#import "Logpackermobilesdk/Logpackermobilesdk.h"
```

Note: Go allows you to build not only SDK but also to compile the application to apk/ipa file from main.go file only without native mobile Ul. 

##Conclusion
Everybody understands that separate commands development for every mobile platform – is not a cheap and easy task. But it is essential for creation a high-quality product at this time. Our task we did in terms of cross platform development and used all its advantages:
* Minimal development resources.
* High development speed.
* Simple decision-support in the future.
