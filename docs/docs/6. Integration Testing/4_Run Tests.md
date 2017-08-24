# Run Tests

## Prerequisites
Please be sure that docker is running.

## 1) Run tests automatically (like jenkins does it)
1. Open a terminal (e.g. IntelliJ or Terminal) and navigate to the "akl/integration-test" folder
2. run the `./integrationtest.sh` shell script to run all tests automatically. This is like jenkins run the tests.

## 2) Run a local test environment with docker and IntelliJ (to debug and write your tests local)
1. Open a terminal (e.g. in IntelliJ) and navigate to the "akl/integration-test" folder
2. run the `./integrationtest-dev.sh` shell script to bring up a local docker environment (flamingo, Keycloak, Selenium..). 
3. After executing the script some docker-compose files will be executed. It could take a while to download the images and start the containers. Please be patient.
4. Now you can access e.g. flamingo or Selenium on your local machine. Open a browser and open e.g. `flamingo:3210`. You should see the AKL shop page.
5. Set the "Arguments:" in your run configuration to one of the options in 2.1)
6. Run your test in Intellij

### 2.1) Options/Configurations to run tests
Is it possible to run the tests in different ways. 

***IntelliJ:*** Please adjust your run configuration to one of the following options (Run -> Edit Configurations... -> Your created gradle run configurations -> Arguments:) \

1. Run your test with a local chrome webdriver against your local docker environment \
`-PtestTarget=compose -Pgeb.env=chrome`

2. Run your test with the selenium docker container against your local docker environment \
`-PtestTarget=compose -Pgeb.env=local` \
(use a vnc client to connect to your selenium container -> e.g. mac os -  "Gehe zu" -> "Mit Server verbinden" -> `vnc://localhost:5900`(password: seceret)

## 3) Run a local test environment with docker and your Terminal
1. Execute one shell script like in 1) or 2) described
2. Open a Terminal an navigate to the akl/integration-test folder. Run
e.g. `gradle -PtestTarget=compose -Pgeb.env=chrome` or `-PtestTarget=compose -Pgeb.env=local`

## 4) Shell script details: 
### 4.1) Complete automated (e.g. for Jenkins)
`./integrationtest.sh`

### 4.2) Complete automated, with VNC and local access

`./integrationtest-dev.sh`

This can be used with Intellij.
