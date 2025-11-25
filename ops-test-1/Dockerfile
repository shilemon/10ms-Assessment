FROM python:3.11-slim

WORKDIR /src

COPY app .
RUN pip install -r requirements.txt

EXPOSE 80

CMD ["python", "main.py"]

