FROM alpine:3.4

WORKDIR /apl-loc-deploy
#RUN mkdir -p /apl-loc-deploy/interviews
COPY artifacts/apl-loc-deploy-linux.tgz .
RUN tar xzvf ./apl-loc-deploy-linux.tgz
ENV APL_LAUNCER="k8s"

CMD ["apl-loc-deploy","--answer-file=silent.json","--silent=true"]