package slt 

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"
	
	"github.com/brunetto/goutils/debug"
)

func DockerBuild (targetHost string) () {
	defer debug.TimeMe(time.Now())
	
	const layout = "2006-01-02-15.04"
	
	var ( 
		envStruct = &EnvStruct{}
		// template to be filled
		baseTmpl *template.Template = template.Must(template.New("baseNameTmpl").Parse(Dockerfile))
		// buffer to write into
		buf bytes.Buffer
		outFile *os.File
		err error
		build *exec.Cmd
		cmd []string 
	)
	
	if targetHost = os.Args[1]; targetHost == "longisland" {
		log.Println("Building starlab docker image for longisland")
		envStruct.EnvString = `##############
# longisland #
##############
ENV CUDA_DRIVER 340.46
ENV CUDA_INSTALL http://us.download.nvidia.com/XFree86/Linux-x86_64/${CUDA_DRIVER}/NVIDIA-Linux-x86_64-${CUDA_DRIVER}.run
ENV CUDA_TOOLKIT cuda_6.0.37_linux_64.run
ENV CUDA_TOOLKIT_DOWNLOAD http://developer.download.nvidia.com/compute/cuda/6_0/rel/installers/$CUDA_TOOLKIT
##############
`
	} else if targetHost == "spritz" {
		log.Println("Building starlab docker image for spritz")
		envStruct.EnvString = `##############
#   spritz   #
##############
ENV CUDA_DRIVER 331.113
ENV CUDA_INSTALL http://us.download.nvidia.com/XFree86/Linux-x86_64/${CUDA_DRIVER}/NVIDIA-Linux-x86_64-${CUDA_DRIVER}.run
ENV CUDA_TOOLKIT cuda_5.5.22_linux_64.run
ENV CUDA_TOOLKIT_DOWNLOAD http://developer.download.nvidia.com/compute/cuda/5_5/rel/installers/$CUDA_TOOLKIT
##############
`
	} else if targetHost == "uno" {
		log.Println("Building starlab docker image for uno")
		envStruct.EnvString = `##############
#    uno     #
##############
ENV CUDA_DRIVER 331.38
ENV CUDA_INSTALL http://us.download.nvidia.com/XFree86/Linux-x86_64/${CUDA_DRIVER}/NVIDIA-Linux-x86_64-${CUDA_DRIVER}.run
ENV CUDA_TOOLKIT cuda_5.5.22_linux_64.run
ENV CUDA_TOOLKIT_DOWNLOAD http://developer.download.nvidia.com/compute/cuda/5_5/rel/installers/$CUDA_TOOLKIT
##############
`
	} else {
		log.Fatal("Unknown host: ", targetHost)
	}
	
	// execute the template and check for errors
	if err := baseTmpl.Execute(&buf, envStruct); err != nil {
		log.Println("Error while creating basename in conf.BaseName:", err)
	}

	if outFile, err = os.Create("Dockerfile"); err != nil {
		log.Fatal("Can't create Dockerfile with error: ", err)
	}
	defer outFile.Close()
	
	if _, err = outFile.Write(buf.Bytes()); err != nil {
		log.Fatal("Can't write to Dockerfile with error: ", err)
	}
	buf.Reset() // reset the buffer just to be sure (not necessary now)
	
	log.Println("Wrote Dockerfile")
	
	cmd = []string{"docker", "build", "--no-cache", "--force-rm=true", "--tag=brunetto/starlab-mapelli-" + targetHost + ":" + time.Now().Format(layout), "."}
	
	log.Println("Run docker build on Dockerfile")

	fmt.Println(cmd)
	
	build = exec.Command(cmd[0], cmd[1:]...)
	if build.Stdout = os.Stdout; err != nil {
		log.Fatal("Error connecting STDOUT: ", err)
	}
	if build.Stderr = os.Stderr; err != nil {
		log.Fatal("Error connecting STDERR: ", err)
	}
	
	if err = build.Start(); err != nil {
		log.Fatal("Error starting build: ", err)
	}
	
	if err = build.Wait(); err != nil {
		log.Fatal("Error while waiting for build: ", err)
	}
}

type EnvStruct struct {
	EnvString string
}

var Dockerfile string = `FROM ubuntu:14.04

MAINTAINER brunetto ziosi <brunetto.ziosi@gmail.com>

# Public version of StarLab4.4.4, see http://www.sns.ias.edu/~starlab/

ENV DEBIAN_FRONTEND noninteractive

ENV STARLAB_FOLDER starlabDocker

# Copy StarLab bundle into the image
COPY $STARLAB_FOLDER /$STARLAB_FOLDER

# This has to be set by hand and MUST be the same of the host
{{.EnvString}}

# Update and install minimal and clean up packages
RUN apt-get update --quiet && apt-get install --yes \
--no-install-recommends --no-install-suggests \
build-essential module-init-tools wget libboost-all-dev   \
&& apt-get clean && rm -rf /var/lib/apt/lists/*

# Install CUDA drivers
RUN wget $CUDA_INSTALL -P /tmp --no-verbose \
      && chmod +x /tmp/NVIDIA-Linux-x86_64-${CUDA_DRIVER}.run \
      && /tmp/NVIDIA-Linux-x86_64-${CUDA_DRIVER}.run -s -N --no-kernel-module \
      && rm -rf /tmp/*

# Install CUDA toolkit
RUN wget $CUDA_TOOLKIT_DOWNLOAD && chmod +x $CUDA_TOOLKIT \
&& ./$CUDA_TOOLKIT -toolkit -toolkitpath=/usr/local/cuda-site -silent -override \
&& rm $CUDA_TOOLKIT

# Set env variables
RUN echo "PATH=$PATH:/usr/local/cuda-site/bin" >> .bashrc          \
&& echo "LD_LIBRARY_PATH=/usr/local/cuda-site/lib64" >> .bashrc   \
&& . /.bashrc \
&& ldconfig /usr/local/cuda-site/lib64

# Install StarLab w/ and w/o GPU, w/ and w/o tidal fields
RUN cp -r /$STARLAB_FOLDER/starlab starlab-GPU           \
&&  cp -r /$STARLAB_FOLDER/starlab starlab-no-GPU        \
&&  cp -r /$STARLAB_FOLDER/starlab starlabAS-GPU         \
&&  cp -r /$STARLAB_FOLDER/starlab starlabAS-no-GPU      \
&&  cp -r /$STARLAB_FOLDER/sapporo sapporo

# Tidal field version only has 5 files different, 
# so we can copy them into a copy of the non TF version:

# starlab/src/node/dyn/util/add_tidal.C
# starlab/src/node/dyn/util/dyn_external.C
# starlab/src/node/dyn/util/dyn_io.C
# starlab/src/node/dyn/util/set_com.C
# starlab/src/node/dyn/util/dyn_story.C

RUN cp /$STARLAB_FOLDER/starlabAS/*.C starlabAS-no-GPU/src/node/dyn/util/ \
&&  cp /$STARLAB_FOLDER/starlabAS/*.C starlabAS-GPU/src/node/dyn/util/     \
&&  cp /$STARLAB_FOLDER/starlabAS/dyn.h starlabAS-no-GPU/include/          \
&&  cp /$STARLAB_FOLDER/starlabAS/dyn.h starlabAS-GPU/include/

# Compile sapporo
RUN cd sapporo/ && make && bash compile.sh && cd ../

# With and w/o GPU and w/ and w/o AS tidal fields
RUN cd /starlab-GPU/ && ./configure --with-f77=no && make && make install && cd ../ \
&& mv /starlab-GPU/usr/bin slbin-GPU \
&& cd /starlabAS-GPU/ && ./configure --with-f77=no && make && make install && cd ../ \
&& mv /starlabAS-GPU/usr/bin slbinAS-GPU \
&& cd /starlab-no-GPU/ && ./configure --with-f77=no --with-grape=no && make && make install && cd ../ \
&& mv /starlab-no-GPU/usr/bin slbin-no-GPU \
&& cd /starlabAS-no-GPU/ && ./configure --with-f77=no --with-grape=no && make && make install && cd ../ \
&& mv /starlabAS-no-GPU/usr/bin slbinAS-no-GPU 

# Default command.
ENTRYPOINT ["/bin/bash"]

`
