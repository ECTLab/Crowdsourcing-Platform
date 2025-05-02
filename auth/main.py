import json
import uuid
from fastapi import FastAPI, Header, HTTPException, Request
from pydantic import BaseModel
import redis

r = redis.Redis(host='localhost', port=6379, db=0, decode_responses=True)

app = FastAPI()


class TokenRequest(BaseModel):
    name: str


@app.post("/generate-token")
def generate_token(req: TokenRequest):
    if not req.name:
        raise HTTPException(status_code=400, detail="send name for the requested token")
    token = uuid.uuid4().hex
    data = {
        "name": req.name,
        "nav": 0,
        "crowd": 0
    }
    r.set(token, json.dumps(data))
    return {"token": token}


@app.get("/bill")
def aggregate(api_key: str = Header(..., alias="api-key")):
    token_data = r.get(api_key)
    if not token_data:
        raise HTTPException(status_code=404, detail="Token not found")

    data = json.loads(token_data)

    nav_keys = r.keys(f"{api_key}_nav_*")
    crowd_keys = r.keys(f"{api_key}_crowd_*")

    nav_count = len(nav_keys)
    crowd_count = len(crowd_keys)

    for key in nav_keys + crowd_keys:
        r.delete(key)

    data["nav"] += nav_count
    data["crowd"] += crowd_count

    r.set(api_key, json.dumps(data))

    fee_per_service = 1
    total_nav_fee = data["nav"] * fee_per_service
    total_crowd_fee = data["crowd"] * fee_per_service

    return {
        "nav": total_nav_fee,
        "crowd": total_crowd_fee
    }
