VOTING_ID=
ALTERNATIVE_ID=
ZONES=("us-central1-a" "us-central1-b" "us-central1-c" "us-west1-a" "us-west1-b" "us-west1-c" "us-east1-b" "us-east1-c" "us-east1-d")

deploy_machine(){
  bash -c "gcloud beta compute --project=the-mesh-296114 instances create votatron-$1 --zone=$1 --machine-type=f1-micro --subnet=default --network-tier=PREMIUM --metadata=startup-script=apt\ get\ update$'\n'curl\ -fsSL\ get.docker.com\ \|\ sh$'\n'docker\ pull\ randofelipe/votatron:1.0.2$'\n'docker\ run\ -d\ -p2112:2112\ --restart\ always\ randofelipe/votatron:1.0.2\ --votingID\ $2\ --alternativeID\ $3 --no-restart-on-failure --maintenance-policy=TERMINATE --preemptible --service-account=1062298037843-compute@developer.gserviceaccount.com --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append --min-cpu-platform=Automatic --no-shielded-secure-boot --shielded-vtpm --shielded-integrity-monitoring --reservation-affinity=any"  
}

for zone in ${ZONES[@]};do
    deploy_machine $zone $VOTING_ID $ALTERNATIVE_ID
done
