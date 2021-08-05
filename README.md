**PHPD CODING CHALLENGE PROJECT**

**Models Files**

**Database Schemas**
schema diagram - https://dbdiagram.io/d/610bf67b2ecb310fc3bf3638
schema statement - https://drive.google.com/file/d/10Pc27iOzyytIIZc8WOhxQ3frstjer7fJ/view?usp=sharing


## Technologies Used
Below is a list of technologies used to build this project

-   Golang
-   Postgres 9.5
-   Makefile

## Installation
Follow these steps to set up the app.

Clone the repo:

[Github repo](https://github.com/rossi1/coding-challenge)

Navigate to the project directory:

`cd coding-challenge`

Note, Ensure you have postgres, Go, Makefile installed on local machine.


Environment variables like  `DB_HOST`,  `DB_PORT`, `DB_USER`, `DB_PASSWORD` and `DB_NAME` should be added. Kindy refer to the .env_example to populate the environmental variable


## Running

Make sure your database is configured and you declared your table with the database schema declared above and your environment variable are set.

## Running with Makefile

Run the command below to run the program

`$ make all`

To build the program binary for different os

 `$ make compile`

Yay!!! You are good to go.

Quote:
*Always **code** as if the guy who ends up maintaining your **code** will be a violent psychopath who knows where you live*