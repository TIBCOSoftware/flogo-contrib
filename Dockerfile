FROM scratch

VOLUME  /flogo/flogo-contrib
CMD ["/bin/true"]
COPY README.md /flogo/flogo-contrib/
COPY activity/ /flogo/flogo-contrib/activity
COPY model/ /flogo/flogo-contrib/model
COPY trigger/ /flogo/flogo-contrib/trigger
