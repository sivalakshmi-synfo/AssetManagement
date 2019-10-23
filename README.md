# AssetManagement
instal mysql driver to connect my sql with golang using "_ github.com/go-sql-driver/mysql" package
>>>>> go get -u github.com/go-sql-driver/mysql
________________________________________________________________
Database 
create database AMS

CREATE TABLE `employees` (
   `IdEmployees` int(11) NOT NULL AUTO_INCREMENT,
   `Name` varchar(45) DEFAULT NULL,
   `DOB` date DEFAULT NULL,
   `Email` varchar(45) DEFAULT NULL,
   `Mobile` varchar(40) DEFAULT NULL,
   `Address` varchar(450) DEFAULT NULL,
   PRIMARY KEY (`IdEmployees`)
 ) 
 
INSERT INTO `ams`.`employees`
(
`Name`,
`DOB`,
`Email`,
`Mobile`,
`Address`)
VALUES
(
'JOHN',
30/05/1888,
'JOHN@synfosys.com',
'1234567889',
'US');
