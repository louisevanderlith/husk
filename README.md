# husk
Husk is a Object-Oriented Database

Husk largely forces users to contain business logic close to their required objects, and minimizes entry points and chances for loop holes. ie. when a developer accesses or modifies data outside of the intended scope. This includes many sources; 

* Directly modifying data an external tool. (SSMS, WorkBench, Toad)
* Able to access core "Database" API in higher-level code, rather than the Logic Layer as intended.  (Front-end, Public facing API End-points)

All records are internally sorted by their creation Timestamp, and then their traditional "ID".
This combination is refered to as a Key. 

#Bench History (TestInserts_SampleETL):
Please note these numbers come from our Sample_ETL test, which inserts the same record(16kb) for 20seconds
* 0.1 (One Record, One File) Write: 138rec/s
* 0.2 (BigFile) Write: 509rec/s (x3.6)
* 0.3 (Dump Index only on save) Write: 1463rec/s (x3)
* 0.4 (Better File handling) Write: 1221rec/s (0%)
* 0.5 (Index Refactor, keys are Ptrs, improved read) Write: 2315rec/s (x2)

--MAC 3167rec/s (Unicorn Power)
--WINDOWS 2315/rec/s (Spinning Disk, AMD)
--LINUX 2289rec/s (SSD, Intel i5(2nd))

# Database Engine
* Data-orientation and clustering
* Everything related to an object will always remain nested within that object. 
* ISAM?


#trying to do this... https://en.wikipedia.org/wiki/IBM_Informix_C-ISAM

##--encourage goroutines. Trigger events for AfterCommit()
* implement golang ReaderWriter for Block Storage