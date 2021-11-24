import grpc
import codegen.oomd_pb2
import codegen.oomd_pb2_grpc

class Client(object):
    def __init__(self):
        pass

    def online_get(self, entity_key, feature_names):
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.OnlineGet(codegen.oomd_pb2.OnlineGetRequest(entity_key=entity_key, feature_names=feature_names))
        return response.result.map

if __name__ == "__main__":
    client = Client()
    print(client.online_get("1006", ["state", "credit_score", "account_age_days", "has_2fa_installed", "transaction_count_7d", "transaction_count_30d"]))
