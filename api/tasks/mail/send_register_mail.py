import logging
from typing import Optional

from celery import shared_task
from flask import current_app, render_template

from extensions.ext_mail import mail


@shared_task
def send_register_mail_task(language: str, to: str, code: str):
    """
    Send register mail task.
    :param language: language
    :param to: email
    :param code: code
    :return:
    """
    try:
        if language == 'zh-Hans':
            subject = "注册验证码"
            html = render_template(
                'register_code_zh.html',
                code=code
            )
        else:
            subject = "Registration Verification Code"
            html = render_template(
                'register_code_en.html',
                code=code
            )

        mail.send(
            subject=subject,
            to=to,
            html=html
        )
    except Exception as e:
        logging.exception(f"Send register mail error: {e}")
