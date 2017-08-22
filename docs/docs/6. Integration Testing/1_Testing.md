# Testing

## 1.) Testing with spock and gradle
### 1.1) Prerequisites
Please check:
1. if IntelliJ is installed on your os
2. if java sdk is installed on your os
3. if groovy is installed on your os, otherwise please download and install the latest groovy version (http://groovy-lang.org/download.html).
4. if groovy plugin is installed under File -> Settings -> Plugins in IntelliJ. if not, go further with "Install JetBrains plugin" and pick it up from list

#### 1.1.1) Prepare your project (IntelliJ)
1. Open the flamingo project in IntelliJ
2. Select your project folder ***flamingo*** and go to (***mac os***) File -> Project structure... 
3. Select "Modules" -> "add" -> "Import Module" 
4. Select the "akl/integration-test" folder and press open 
5. Select Gradle and press "Next" 
6. Check that only "Create separate module per source set" is marked (checkboxes) 
7. Check that a "Gradle JVM:" is selected and press "finish" 
8. Now you should see the "integration-test" project underneath the "flamingo" project 

#### 1.1.2) Hosts File
Please add `127.0.0.1 flamingo` to your hosts file.

### 1.2) Strcuture
#### 1.2.1) Location Tests
All integration-tests are stored in ***src -> test -> groovy -> com.aoe.om3.akl -> specs.gui*** \
All new tests should be placed in this folder.

#### 1.2.2) Location Page Objects
All Page Objects are stored in ***src -> test -> groovy -> com.aoe.om3.akl -> pageObjects*** \
All new Page Objects should be placed in this folder.

#### 1.2.3) Location Modules
All Modules are stored in ***src -> test -> groovy -> com.aoe.om3.akl -> modules*** \
All new Modules should be placed in this folder.

### 1.3) Write a test
Please see [Groovy and Spock](https://extranet.aoe.com/confluence/display/aoetestteam/Spock+and+Gradle) for an Groovy and Spock example
