# Lockwood Presence Service

## Assumptions
* Ids are positive ints
* Friendship is a two way street...
* * If User-1 has User-2 as a friend then User-2 MUST have User-1 as a friend, let's call this `symmetrical` friendship
* * If the friendship is not `symmetrical` (if it is `asymmetrical`) then online/offline notifications will not happened
* * It is not within scope to ensure/manage this symmetry, however we should be able to guard against `asymmetrical` friends recieving online/offline notification
* * in the case of `asymmetrical` friendship the unrecipricating user will appear to be offline
