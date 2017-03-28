FROM connecthq/scratch-ssl
ADD app /app
EXPOSE 8080
CMD ["/app"]