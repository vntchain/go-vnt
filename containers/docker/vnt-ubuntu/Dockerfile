FROM ubuntu:xenial

ENV PATH=/usr/lib/go-1.9/bin:$PATH

ADD ./clang.sh /

ADD ./genesis.json.templet /

ADD ./run.sh  /

RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends lsb software-properties-common golang-1.9 git make gcc libc-dev ca-certificates wget && \
  /bin/bash ./clang.sh 5.0 && \
  ln -s /usr/lib/llvm-5.0/lib/libclang*so /usr/lib/ && \
  git clone --depth 1 https://github.com/vntchain/go-vnt -b docker && \
  (cd go-vnt && make all) && \
  cp go-vnt/build/bin/* / && \
  apt-get remove -y golang-1.9 git make gcc libc-dev && apt autoremove -y && apt-get clean && \
  rm -rf /go-vnt


EXPOSE 8555
EXPOSE 30311

EXPOSE 8556
EXPOSE 30312

EXPOSE 8557
EXPOSE 30313

EXPOSE 8558
EXPOSE 30314

EXPOSE 8559
EXPOSE 30315

EXPOSE 8560
EXPOSE 30316

EXPOSE 8561
EXPOSE 30317

EXPOSE 8562
EXPOSE 30318

EXPOSE 8563
EXPOSE 30319

EXPOSE 8564
EXPOSE 30320

EXPOSE 8565
EXPOSE 30321

EXPOSE 8566
EXPOSE 30322

EXPOSE 8567
EXPOSE 30323

EXPOSE 8568
EXPOSE 30324

EXPOSE 8569
EXPOSE 30325

EXPOSE 8570
EXPOSE 30326

EXPOSE 8571
EXPOSE 30327

EXPOSE 8572
EXPOSE 30328

EXPOSE 8573
EXPOSE 30329

EXPOSE 8574
EXPOSE 30330

EXPOSE 8575
EXPOSE 30331

EXPOSE 8576
EXPOSE 30332

EXPOSE 8577
EXPOSE 30333

EXPOSE 8578
EXPOSE 30334


ENTRYPOINT ["bash","/run.sh"]
