css:
	./tailwindcss -i static/styles/input.css -o static/styles/output.css --minify

css-watch:
	./tailwindcss -i static/styles/input.css -o static/styles/output.css --watch

run:
	go run .