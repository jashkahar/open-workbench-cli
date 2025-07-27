from fastapi import FastAPI

app = FastAPI(
    title="{{.ProjectName}}",
    description="A new project scaffolded by Open Workbench CLI.",
    version="0.1.0",
)

@app.get("/")
def read_root():
    return {"message": "Welcome to {{.ProjectName}}!", "owner": "{{.Owner}}"}