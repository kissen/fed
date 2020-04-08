-> removing likes; this involves updating webvocab to know what was
   liked and sending tombstones/deletes/whatever

-> database transactions; implement it in FedStorage and review all
   functions dealing with storage, in particular in ap/; for example
   updateActor is racy w/o transactions

-> once poc is completed: serious refactor

-> support running w/o nginx and https; currently we expect al IRIs
   of our instance to be reachable via https

-> create complete list of reserved usernames
