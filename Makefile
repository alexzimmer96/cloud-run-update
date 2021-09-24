functionName = update-cloud-run
region = europe-west1
serviceAccount = SOME-SERVICE-ACCOUNT

endpoint = ""
registry = ""
project = ""

deploy:
	gcloud functions deploy $(functionName) \
		--region=$(region) \
		--memory=128MB \
		--entry-point Update \
		--trigger-topic gcr \
		--runtime go116 \
		--service-account=$(serviceAccount) \
		--max-instances 1 \
		--set-env-vars ENDPOINT=$(endpoint),REGISTRY=$(registry),PROJECT=$(project)