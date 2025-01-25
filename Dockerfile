FROM archlinux
WORKDIR /app
RUN pacman -Syu --noconfirm
RUN pacman -S go prettier --noconfirm
COPY . .
CMD ["go", "run", "."]
