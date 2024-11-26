FROM archlinux
WORKDIR /app
COPY . .
RUN pacman -Sy go --needed --noconfirm
CMD [ "go", "run", "." ]
