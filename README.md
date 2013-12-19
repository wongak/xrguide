# X Rebirth Guide

## The Tools

### Text Importing

The command `import_t` will import the text from the text directory.

**Usage:**

Calling `import_t` with the -h flag will give you a list of understood parameters.

Here is a list of the currently implemented parameters and a short explanation:

    -dbt The database type. Currently "sqlite3" and "mysql" are supported. But many more backends can be added easily.
	    Default: "sqlite3"
	    Example: 
			-dbt mysql 
			-dbt "mysql"
			-dbt=sqlite3
			-dbt="mysql"
		
	-dsn The DSN (Data Source Name) is specific to the selected database type.
	    Default: "xrguide.db" - This is for the sqlite3 backend and will use the xrguide.db file.
		Example:
			-dsn german.db
			-dsn root:root@tcp(localhost:3306)/xrguide?charset=utf8mb4,utf8
				This is a MySQL DSN. 
				You can find the specifics here: https://github.com/go-sql-driver/mysql#dsn-data-source-name
				
	-l The language id to be imported. If 0, all languages will be imported.
		Example:
			-l 44
				For English only.
			-l 49
				For German only.
	
	-p The page id to be imported. If 0, all pages will be imported.
		Example: 
			-p 1001
				Will only import the page 1001 (Interface)
				
	-r Whether to recreate the tables. All the old data will be lost!
		Example:
			-r=1
			-r=0
		Please note that the equals sign "=" is required for boolean flags.
		
	-t The path to the text files.
		Default: "." Current directory
		Example:
			-t ../resources/t
			
	-v Verbose output. This will be significantly slower.
		Example:
			-v=1
			
	-w Worker count. Modifies the concurrent pages to be imported. Shouldn't be too high, since the bottleneck of the importing will be the disk I/O.
		Default: 10
		Example:
			-t 50
