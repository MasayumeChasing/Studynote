FROM registry-cbu.huawei.com/euleros_docker/euleros_x86_64:2.9.5

ENV app=mps

RUN useradd -u 1001 $app && \
    useradd -u 1337 istio  && \
	rpm -qa|grep seccomponent || rpm -Uvh http://obs.cn-north-5.myhuaweicloud.com/fuxi-yum/x86_64/test/Packages/seccomponent-1.0.1-1.5e.release.x86_64.rpm && \
    mkdir -p /opt/{$app/conf,certs,scc/{pk,stb},holo/amah_agent,log/$app/{run,scc,operation}} && \
    chown -R $app:$app /opt && \
    chmod -R 700 /opt

COPY --chown=1001:1001 . /opt/$app

ENV LD_LIBRARY_PATH=/opt/mps/libmps_common.so:$LD_LIBRARY_PATH

USER $app

WORKDIR /opt/$app

CMD ["./entrypoint.sh"]


