import secrets
from datetime import UTC, datetime
import random

from flask import request
from flask_restful import Resource, reqparse
from sqlalchemy import select
from sqlalchemy.orm import Session

from controllers.console import api
from controllers.console.error import (
    AccountInFreezeError,
    EmailCodeError,
    EmailPasswordResetLimitError,
    EmailSendIpLimitError,
    InvalidEmailError,
    InvalidTokenError,
    PasswordMismatchError,
)
from controllers.console.wraps import setup_required
from extensions.ext_database import db
from libs.helper import extract_remote_ip
from libs.password import hash_password, valid_password
from models.account import Account
from services.account_service import AccountService, TenantService
from services.errors.account import AccountRegisterError
from services.errors.workspace import WorkSpaceNotAllowedCreateError
from services.feature_service import FeatureService
from tasks.mail.send_register_mail import send_register_mail_task
from libs.token_manager import TokenManager
from config import languages


class RegisterSendEmailApi(Resource):
    @setup_required
    def post(self):
        parser = reqparse.RequestParser()
        parser.add_argument("email", type=str, required=True, location="json")
        parser.add_argument("language", type=str, required=False, location="json")
        args = parser.parse_args()

        ip_address = extract_remote_ip(request)
        if AccountService.is_email_send_ip_limit(ip_address):
            raise EmailSendIpLimitError()

        if args["language"] is not None and args["language"] == "zh-Hans":
            language = "zh-Hans"
        else:
            language = "en-US"

        with Session(db.engine) as session:
            account = session.execute(select(Account).filter_by(email=args["email"])).scalar_one_or_none()
        
        if account is not None:
            return {"result": "fail", "code": "account_already_exists", "message": "Account already exists"}

        # 发送注册验证码邮件
        code = "".join([str(random.randint(0, 9)) for _ in range(6)])
        token = TokenManager.generate_token(
            account=None, email=args["email"], token_type="register", additional_data={"code": code}
        )
        send_register_mail_task.delay(
            language=language,
            to=args["email"],
            code=code,
        )

        return {"result": "success", "data": token}


class RegisterCheckApi(Resource):
    @setup_required
    def post(self):
        parser = reqparse.RequestParser()
        parser.add_argument("email", type=str, required=True, location="json")
        parser.add_argument("code", type=str, required=True, location="json")
        parser.add_argument("token", type=str, required=True, nullable=False, location="json")
        args = parser.parse_args()

        user_email = args["email"]

        is_forgot_password_error_rate_limit = AccountService.is_forgot_password_error_rate_limit(args["email"])
        if is_forgot_password_error_rate_limit:
            raise EmailPasswordResetLimitError()

        token_data = AccountService.get_reset_password_data(args["token"])
        if token_data is None:
            raise InvalidTokenError()

        if user_email != token_data.get("email"):
            raise InvalidEmailError()

        if args["code"] != token_data.get("code"):
            AccountService.add_forgot_password_error_rate_limit(args["email"])
            raise EmailCodeError()

        AccountService.reset_forgot_password_error_rate_limit(args["email"])
        return {"is_valid": True, "email": token_data.get("email")}


class RegisterCompleteApi(Resource):
    @setup_required
    def post(self):
        parser = reqparse.RequestParser()
        parser.add_argument("token", type=str, required=True, nullable=False, location="json")
        parser.add_argument("name", type=str, required=True, nullable=False, location="json")
        parser.add_argument("password", type=valid_password, required=True, nullable=False, location="json")
        parser.add_argument("password_confirm", type=valid_password, required=True, nullable=False, location="json")
        args = parser.parse_args()

        # 验证密码匹配
        if args["password"] != args["password_confirm"]:
            raise PasswordMismatchError()

        # 验证token并获取数据
        reset_data = AccountService.get_reset_password_data(args["token"])
        if not reset_data:
            raise InvalidTokenError()

        # 撤销token防止重用
        AccountService.revoke_reset_password_token(args["token"])

        email = reset_data.get("email", "")

        # 创建新账号
        try:
            AccountService.create_account_and_tenant(
                email=email,
                name=args["name"],
                password=args["password"],
                interface_language=languages[0],
            )
        except WorkSpaceNotAllowedCreateError:
            pass
        except AccountRegisterError:
            raise AccountInFreezeError()

        return {"result": "success"}


api.add_resource(RegisterSendEmailApi, "/register")
api.add_resource(RegisterCheckApi, "/register/validity")
api.add_resource(RegisterCompleteApi, "/register/complete")
