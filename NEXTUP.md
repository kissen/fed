-> replace w/ generic version

	user, err := retrieveOwner(&iri, fedcontext.From(c).Storage)
	if err != nil {
		return nil, err
	}

	liked = prop.ToCollection(user.Liked)

	id := fediri.LikedIRI(user.Name).URL()
	prop.SetIdOn(liked, id)

	return liked, nil

-> database transactions; implement it in FedStorage and review all
   functions dealing with storage, in particular in ap/; for example
   updateActor is racy w/o transactions

-> once poc is completed: serious refactor
