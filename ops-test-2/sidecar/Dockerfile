FROM python:3.9-slim

WORKDIR /app

COPY proxy.py /app/

RUN pip install flask requests

EXPOSE 8081

CMD ["python3", "proxy.py"]
