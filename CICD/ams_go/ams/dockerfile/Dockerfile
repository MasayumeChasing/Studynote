FROM registry-cbu.huawei.com/euleros_docker/euleros_x86_64:2.9.5

ENV  app=ams

RUN useradd -u 1001 $app && \
    useradd -u 1337 istio  && \
	rpm -Uvh http://obs.cn-north-5.myhuaweicloud.com/fuxi-yum/x86_64/test/Packages/seccomponent-1.0.1-1.5e.release.x86_64.rpm && \
    mkdir -p /opt/{$app/conf,scc/{access,media,user}/{pk,stb},log/$app/{run,scc,operation}} && \
    chown -R $app:$app /opt && \
    chmod -R 700 /opt && \
    mkdir -p /log/holo && \
    chown -R $app:$app /log && \
    chmod -R 700 /log

COPY --chown=1001:1001 . /opt/$app

USER $app

WORKDIR /opt/$app

EXPOSE 7070

CMD ["./entrypoint.sh"]


