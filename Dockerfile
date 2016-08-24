FROM scratch

VOLUME  /flogo/contrib
COPY activity /flogo/contrib/
COPY model /flogo/contrib/
COPY trigger /flogo/contrib/
COPY README.md /flogo/contrib/
