-> removing likes; this involves updating webvocab to know what was
   liked and sending tombstones/deletes/whatever

-> database transactions; implement it in FedStorage and review all
   functions dealing with storage, in particular in ap/; for example
   updateActor is racy w/o transactions

-> once poc is completed: serious refactor
