FROM archlinux
WORKDIR /app
COPY . .
RUN pacman -Sy prettier go --needed --noconfirm
CMD [ "go", "run", "." ]
