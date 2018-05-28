apl-loc-deploy is a tool to create kubernetes clusters externally (not via appLariat UI).  It is a Q&A based program and wil

Requirements
------------
    - docker must be installed where the program is executed
    - appLariat credential id for your desired cloud provider
    - If importing an existing cluster, fqdn and credentials are needed
    - Setting of two variables, APL_API and APL_API_KEY.  APL_API_KEY is generated via the UI.
        export APL_API="https://api.applariat.io/v1"
        export APL_API_KEY="your api key"

Usage
------
    - Run the program apl-loc-deploy and answer the questions
    - If there is a default answer shown in (), and you want to accept it, simply press enter.
    - Use -id=identifier to use previously used interview questions and answers
      apl-loc-deploy -id=applariat-11-15  this will use all files with *_applariat-11-15.json for the interview

Interview Questions
--------------------
    - Questions the program asks are in json files in the interviews directory
        - general.json - General questions that are asked for any type of cluster
        - *.json - The other json files are questions per provider or other functionality
        - Saved files will also go in here with the name of the cluster: i.e - general_cluster1.json



