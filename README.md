# husk
Husk is a Object-Oriented Database

Husk largely forces users to contain business logic close to their required objects, and minimizes entry points and chances for loop holes.

All records are internally sorted by their creation Timestamp, and then their traditional "ID".
This combination is refered to as a Key. 

Bench History (TestInserts_SampleETL):
* 0.1 (One Record, One File) Write: 138rec/s
* 0.2 (BigFile) Write: 509rec/s (x3.6)
* 0.3 (Dump Index only on save) Write: 1463rec/s (x3)
* 0.4 (Better File handling) Write: 1221rec/s (0%)
* 0.5 (Index Refactor, keys are Ptrs, improved read) Write: 2315rec/s (x2)

# Database Engine
* Data-orientation and clustering
Everything related to an object will always remain nested within that object. 

#trying to do this... https://en.wikipedia.org/wiki/IBM_Informix_C-ISAM
#with some help from this... https://www.codeproject.com/articles/1029838/build-your-own-database


##--encourage goroutines. Trigger events for AfterCommit()
* implement golang ReaderWriter for Block Storage