FROM registry.xiaoyou.host/library/python:3.9.16
WORKDIR /code
COPY . .
RUN  pip3 install -r requirements.txt -i https://nexus.xiaoyou.host/repository/pip-hub/simple
EXPOSE 9000
CMD ["python3","main.py"]