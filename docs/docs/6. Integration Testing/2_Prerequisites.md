# Prerequisites

## 1) Prerequisites
Please check:
1. if IntelliJ is installed on your os
2. if java sdk is installed on your os
3. if groovy is installed on your os, otherwise please download and install the latest groovy version (http://groovy-lang.org/download.html).
4. if groovy plugin is installed under File -> Settings -> Plugins in IntelliJ. if not, go further with "Install JetBrains plugin" and pick it up from list
5. Prepare your project
6. Creata a gradle run configuration
7. Add an entry to your hosts file

## 1.1) Prepare your project (IntelliJ)
1. Open the flamingo project in IntelliJ
2. Select your project folder ***flamingo*** and go to (***mac os***) File -> Project structure... 
3. Select "Modules" -> "add" -> "Import Module" 
4. Select the "akl/integration-test" folder and press open 
5. Select Gradle and press "Next" 
6. Check that only "Create separate module per source set" is marked (checkboxes) 
7. Check that a "Gradle JVM:" is selected and press "finish" 
8. Now you should see the "integration-test" project underneath the "flamingo" project
9. Select your new inetragtion-test project and go to "Run" -> "Edit configurations..."
10. In the new windows select gradle and press "plus"
11. In the right window choose "integration-test" folder as "Gradle project:"
12. In the right window put "test" in "Tasks:"
13. In the right window put "-PtestTarget=compose -Pgeb.env=chrome" in "Arguments" (without the quotes) ***If you dont see the "Arguments" filed, please restart IntelliJ***
14. Click on "Apply" and "Ok"
15. Open a Terminal (e.g. in Intellij) and navigate to ***akl/integration-test***
16. Execute `docker login docker-om3-akl.aoe.com` and put in your credentials (e.g. max.mustermann and your password)



## 1.2) Add entries to your hosts file
Please add `127.0.0.1 flamingo` to your hosts file. \
Please add `127.0.0.1 keycloak` to your hosts file. 
