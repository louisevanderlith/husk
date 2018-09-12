# husk
Husk is a Object-Oriented Database

Husk largely forces users to contain business logic close to their required objects, and minimizes entry points and chances for loop holes.

All records are internally sorted by their creation Timestamp, and then their traditional "ID".
This combination is refered to as a Key. 

Bench History (TestInserts_SampleETL):
0.1 (One Record, One File) Write: 138rec/s
0.2 (BigFile) Write: 509rec/s (x3.6)