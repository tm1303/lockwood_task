# Lockwood Presence Service
 
## Build
 
```
docker build -f Dockerfile.server --tag lockwood_task .
docker build -f Dockerfile.client --tag lockwood_task_client .
```
 
## Run
 
Run one instance of the server first
```
docker run --rm -it -p 13131:13131/udp lockwood_task
```
Run as many clients as you like
```
docker run --rm -it --network host lockwood_task_client
```

## Run from DockerHub (not recommended)

These images are big and unoptimised (and who knows what virus's I've packaged up) so not recommended but if you want to run direct from dockerhub...
```
docker run --rm -it -p 13131:13131/udp tm1303/lockwood:v1
docker run --rm -it --network host tm1303/lockwood:v1_client
```
 
## User Guide
I would recommend opening 4 or more terminal windows to get a feel for how this works. In your first terminal run the server as above and you should see `Presence Server listening for user logins`. Later on there will be some debug logging in this window but I'd recommend using it as "clue" to the code rather than trying to understand it directly.
 
In another window run a client as above, you'll be asked for a user-id (an int) with which to login, when you hit return you'll be asked for a friend-id (again an int) you can have many friends (keep track of a few inputs for later), enter them one int at a time and press return after each. Once you have entered a few then hit return again on an empty input and you will be connected...
```
Beginning transmission...
Awaiting notifications...
```
Back in your server instance you'll see `User Connecting: 1` and then the service will periodically refresh your friend's "online status" and print the outcome to the server terminal. You won't see anything in the client yet, we don't want to bombard them.
 
In the next terminal run another client as before, pick a user-id which you specified in your first user's friends list and remember to include your first user in this new user's friend list, as all friend-status must be symmetrical to share online/offline updates. When you complete this login you should see as before but with an addition
```
Beginning transmission...
Awaiting notifications...
> Your friend is ONLINE! (UserId: 1)
```
and back in your first client you should see something similar `> Your friend is ONLINE! (UserId: 2)`. Add a few more terminals. Try adding a few asymmetrical friendships, neither side will receive a notification but you will see the failed check in the server output
```
1 did NOT verify 77
3 did NOT verify 77
```
If you close a client there will be no immediate indication. This is all UDP based so the online status is based on the client pinging the server every few seconds, if it fails to ping the server will "time-out" the user, remove them from the master list of users (which will appear in the server output), and when the next set of status refreshes are asynchronously processed they will notify back to their client output...
```
Your friend is OFFLINE! (UserId: 1)
```
This typically takes 10 seconds or so, but the time-out, the refresh, and the ping can all be tuned more or less frequently (except they're a bit hardcoded at present).
 
In the server terminal you'll see a whole mess at this point. Interestingly if you kill the server and restart it quick enough the clients will carry on unaware.
 
## Bug-like things
* If you kill a client and reconnect with the same id very quickly things break
* You can log in with the same id in two clients, resulting in weird but explainable behaviour
* Validation on the client inputs is not 100%
* in udp_server.go there is a bit of a design failure, nothing terminal though :/
 
## Assumptions
* Ids are positive ints
* I've built this as a UDP client/server, I've hinted at a swappable approach with a few references to TCP but didn't have time to commit to the idea
* Friendship is a two way street...
  * If User-1 has User-2 as a friend then User-2 MUST have User-1 as a friend, let's call this `symmetrical` friendship
  * If the friendship is not `symmetrical` (if it is `asymmetrical`) then online/offline notifications will not happened
  * It is not within scope to ensure/manage this symmetry, however we should be able to guard against `asymmetrical` friends receiving online/offline notification
  * in the case of `asymmetrical` friendship the reciprocating user will appear to be offline
 
## Questions
1. What changes if we switch TCP to UDP?
* TCP requires us to think about ongoing connections rather than just transmitting data between two addresses. TCP is more conncerned with lost data than UDP but UDP is less overhead. 
2. How would you detect the user (connection) is still online?
* My UDP client has to keep reminding the server it is still online, if this was TCP they'd "hold hands" for the lifetime of the connection so no need to remind the server.
3. What happens if the user has a lot of friends?
* I've tried to design in such a way that the immediate behaviour isn't affected by lots of friends but one pottential issue I haven't covered off is if the periodic "refresh" takes longer than the refresh period. The other concern is mutex congestion, as we increase the locking for more users we have a single small mutex to squeeze sync access via, this will bottleneck.
4. How design of your application will change if there will be more than one instance of the server
* My main goal in this case would be to either move the in-mem store out to maybe redis (which would be nice for handling UDP timeouts) or if we stick with in mem we'd need a mechanism for eventual consitancy.
 

