import Link from 'next/link'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useRouter, useSearchParams } from 'next/navigation'
import { useContext } from 'use-context-selector'
import Button from '@/app/components/base/button'
import Toast from '@/app/components/base/toast'
import { emailRegex } from '@/config'
import { login } from '@/service/common'
import Input from '@/app/components/base/input'
import I18NContext from '@/context/i18n'
import { noop } from 'lodash-es'

type MailAndPasswordAuthProps = {
  isInvite: boolean
  isEmailSetup: boolean
  allowRegistration: boolean
}

const passwordRegex = /^(?=.*[a-zA-Z])(?=.*\d).{8,}$/

export default function MailAndPasswordAuth({ isInvite, isEmailSetup, allowRegistration }: MailAndPasswordAuthProps) {
  const { t } = useTranslation()
  const { locale } = useContext(I18NContext)
  const router = useRouter()
  const searchParams = useSearchParams()
  const [showPassword, setShowPassword] = useState(false)
  const emailFromLink = decodeURIComponent(searchParams.get('email') || '')
  const [email, setEmail] = useState(emailFromLink)
  const [password, setPassword] = useState('')

  const [isLoading, setIsLoading] = useState(false)
  const handleEmailPasswordLogin = async () => {
    if (!email) {
      Toast.notify({ type: 'error', message: t('login.error.emailEmpty') })
      return
    }
    if (!emailRegex.test(email)) {
      Toast.notify({
        type: 'error',
        message: t('login.error.emailInValid'),
      })
      return
    }
    if (!password?.trim()) {
      Toast.notify({ type: 'error', message: t('login.error.passwordEmpty') })
      return
    }
    if (!passwordRegex.test(password)) {
      Toast.notify({
        type: 'error',
        message: t('login.error.passwordInvalid'),
      })
      return
    }
    try {
      setIsLoading(true)
      const loginData: Record<string, any> = {
        email,
        password,
        language: locale,
        remember_me: true,
      }
      if (isInvite)
        loginData.invite_token = decodeURIComponent(searchParams.get('invite_token') as string)
      const res = await login({
        url: '/login',
        body: loginData,
      })
      if (res.result === 'success') {
        if (isInvite) {
          router.replace(`/signin/invite-settings?${searchParams.toString()}`)
        }
        else {
          localStorage.setItem('console_token', res.data.access_token)
          localStorage.setItem('refresh_token', res.data.refresh_token)
          // Extend Begin  ----------------
          // 如果本地浏览器缓存数据存在重定向url，则跳转到重定向url
          if (localStorage.getItem('redirect_url')) {
            const redirectUrl = localStorage.getItem('redirect_url')
            localStorage.removeItem('redirect_url')
            router.replace(redirectUrl as string)
            return
          }
          router.replace('/explore/apps-center-extend')
          // Extend End  ----------------
        }
      }
      else if (res.code === 'account_not_found') {
        // 始终允许注册，不再检查 allowRegistration
        console.log('账号不存在，准备注册:', res);

        // 直接跳转到注册页面
        const params = new URLSearchParams()
        params.append('email', encodeURIComponent(email))
        router.replace(`/register?${params.toString()}`)

        Toast.notify({
          type: 'info',
          message: '账号不存在，正在跳转到注册页面',
        });
      }
      else {
        // 显示详细的错误信息
        console.log('登录失败，完整响应:', res);
        Toast.notify({
          type: 'error',
          message: `登录失败: ${res.data || '未知错误'} (代码: ${res.code || '无代码'})`,
        })
      }
    }

    finally {
      setIsLoading(false)
    }
  }

  return <form onSubmit={noop}>
    <div className='mb-3'>
      <label htmlFor="email" className="system-md-semibold my-2 text-text-secondary">
        {t('login.email')}
      </label>
      <div className="mt-1">
        <Input
          value={email}
          onChange={e => setEmail(e.target.value)}
          disabled={isInvite}
          id="email"
          type="email"
          autoComplete="email"
          placeholder={t('login.emailPlaceholder') || ''}
          tabIndex={1}
        />
      </div>
    </div>

    <div className='mb-3'>
      <label htmlFor="password" className="my-2 flex items-center justify-between">
        <span className='system-md-semibold text-text-secondary'>{t('login.password')}</span>
        <Link
          href={`/reset-password?${searchParams.toString()}`}
          className={`system-xs-regular ${isEmailSetup ? 'text-components-button-secondary-accent-text' : 'pointer-events-none text-components-button-secondary-accent-text-disabled'}`}
          tabIndex={isEmailSetup ? 0 : -1}
          aria-disabled={!isEmailSetup}
        >
          {t('login.forget')}
        </Link>
      </label>
      <div className="relative mt-1">
        <Input
          id="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter')
              handleEmailPasswordLogin()
          }}
          type={showPassword ? 'text' : 'password'}
          autoComplete="current-password"
          placeholder={t('login.passwordPlaceholder') || ''}
          tabIndex={2}
        />
        <div className="absolute inset-y-0 right-0 flex items-center">
          <Button
            type="button"
            variant='ghost'
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? '👀' : '😝'}
          </Button>
        </div>
      </div>
    </div>

    <div className='mb-4'>
      <Button
        tabIndex={2}
        variant='primary'
        onClick={handleEmailPasswordLogin}
        disabled={isLoading || !email || !password}
        className="w-full"
      >{t('login.signBtn')}</Button>
    </div>

    {/* 添加注册入口 */}
    <div className="relative mb-4">
      <div className="absolute inset-0 flex items-center" aria-hidden="true">
        <div className='h-px w-full bg-gradient-to-r from-background-gradient-mask-transparent via-divider-regular to-background-gradient-mask-transparent'></div>
      </div>
      <div className="relative flex justify-center">
        <span className="system-xs-medium-uppercase px-2 text-text-tertiary">{t('login.or')}</span>
      </div>
    </div>

    <div className='mb-2'>
      <Button
        tabIndex={3}
        variant='secondary'
        onClick={() => {
          const params = new URLSearchParams()
          params.append('email', encodeURIComponent(email))
          router.push(`/register?${params.toString()}`)
        }}
        className="w-full"
      >{t('login.registerNewAccount') || '注册新账号'}</Button>
    </div>
  </form>
}
