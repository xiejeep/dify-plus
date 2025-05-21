'use client'
import Link from 'next/link'
import { useCallback, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useRouter, useSearchParams } from 'next/navigation'
import cn from 'classnames'
import { RiArrowLeftLine, RiCheckboxCircleFill, RiEyeLine, RiEyeOffLine, RiUserLine } from '@remixicon/react'
import { useCountDown } from 'ahooks'
import Button from '@/app/components/base/button'
import { completeRegistration } from '@/service/common'
import Toast from '@/app/components/base/toast'
import Input from '@/app/components/base/input'

const validPassword = /^(?=.*[a-zA-Z])(?=.*\d).{8,}$/
const AUTO_REDIRECT_TIME = 3000

const RegisterSetInfoPage = () => {
  const { t } = useTranslation()
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = decodeURIComponent(searchParams.get('token') || '')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [showSuccess, setShowSuccess] = useState(false)
  const [leftTime, setLeftTime] = useState<number>()
  const [countdown] = useCountDown({
    targetDate: leftTime,
    onEnd: () => {
      router.replace('/signin')
    },
  })

  const getSignInUrl = () => {
    return '/signin'
  }

  const showErrorMessage = useCallback((message: string) => {
    Toast.notify({
      type: 'error',
      message,
    })
  }, [])

  const valid = useCallback(() => {
    if (!name.trim()) {
      showErrorMessage(t('login.error.nameEmpty') || '请输入姓名')
      return false
    }
    if (!password.trim()) {
      showErrorMessage(t('login.error.passwordEmpty'))
      return false
    }
    if (!validPassword.test(password)) {
      showErrorMessage(t('login.error.passwordInvalid'))
      return false
    }
    if (password !== confirmPassword) {
      showErrorMessage(t('common.account.notEqual'))
      return false
    }
    return true
  }, [name, password, confirmPassword, showErrorMessage, t])

  const handleRegister = useCallback(async () => {
    if (!valid())
      return
    try {
      setIsLoading(true)
      await completeRegistration({
        token,
        name,
        password,
        password_confirm: confirmPassword,
      })
      setShowSuccess(true)
      setLeftTime(AUTO_REDIRECT_TIME)
    }
    catch (error) {
      console.error(error)
    }
    finally {
      setIsLoading(false)
    }
  }, [name, password, token, valid, confirmPassword])

  return (
    <div className='flex min-h-screen flex-col bg-white'>
      <div className='flex flex-1 flex-col items-center justify-center'>
        <div className='w-full max-w-[400px] px-4'>
          {!showSuccess && (
            <>
              <div className='mb-8 flex items-center justify-between'>
                <Link href='/register/check-code' className='group flex items-center text-sm font-normal text-gray-500 hover:text-gray-900'>
                  <RiArrowLeftLine className='mr-1 h-4 w-4 text-gray-500 group-hover:text-gray-900' />
                  {t('login.back')}
                </Link>
              </div>
              <div className='mb-6 flex items-center'>
                <div className='flex h-10 w-10 items-center justify-center rounded-lg border border-gray-100 shadow-xs'>
                  <RiUserLine className='h-5 w-5 text-primary-600' />
                </div>
                <div className='ml-3'>
                  <div className='text-xl font-semibold text-gray-900'>{t('login.setInfo') || '设置账号信息'}</div>
                  <div className='text-sm text-gray-500'>{t('login.setInfoTip') || '请设置您的账号信息'}</div>
                </div>
              </div>
              <div className='mb-4'>
                <label htmlFor="name" className='system-md-semibold my-2 text-text-secondary'>{t('login.name') || '姓名'}</label>
                <div className='mt-1'>
                  <Input id='name' type="text" disabled={isLoading} value={name} placeholder={t('login.namePlaceholder') || '请输入您的姓名' as string} onChange={e => setName(e.target.value)} />
                </div>
              </div>
              <div className='mb-4'>
                <label htmlFor="password" className='system-md-semibold my-2 text-text-secondary'>{t('login.newPassword')}</label>
                <div className='relative mt-1'>
                  <Input
                    id='password'
                    type={showPassword ? 'text' : 'password'}
                    disabled={isLoading}
                    value={password}
                    placeholder={t('login.passwordPlaceholder') as string}
                    onChange={e => setPassword(e.target.value)}
                  />
                  <div className='absolute inset-y-0 right-0 flex items-center'>
                    <Button
                      type='button'
                      variant='ghost'
                      onClick={() => setShowPassword(!showPassword)}
                      className='px-3'
                    >
                      {showPassword ? <RiEyeOffLine className='h-5 w-5 text-gray-500' /> : <RiEyeLine className='h-5 w-5 text-gray-500' />}
                    </Button>
                  </div>
                </div>
                <div className='mt-1 text-xs text-gray-500'>{t('login.passwordRequirement')}</div>
              </div>
              <div className='mb-6'>
                <label htmlFor="confirmPassword" className='system-md-semibold my-2 text-text-secondary'>{t('login.confirmPassword')}</label>
                <div className='relative mt-1'>
                  <Input
                    id='confirmPassword'
                    type={showConfirmPassword ? 'text' : 'password'}
                    disabled={isLoading}
                    value={confirmPassword}
                    placeholder={t('login.confirmPasswordPlaceholder') as string}
                    onChange={e => setConfirmPassword(e.target.value)}
                  />
                  <div className='absolute inset-y-0 right-0 flex items-center'>
                    <Button
                      type='button'
                      variant='ghost'
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      className='px-3'
                    >
                      {showConfirmPassword ? <RiEyeOffLine className='h-5 w-5 text-gray-500' /> : <RiEyeLine className='h-5 w-5 text-gray-500' />}
                    </Button>
                  </div>
                </div>
              </div>
              <div className='mb-4'>
                <Button loading={isLoading} disabled={isLoading} variant='primary' className='w-full' onClick={handleRegister}>{t('login.register') || '注册'}</Button>
              </div>
            </>
          )}

          {showSuccess && (
            <div className="flex flex-col md:w-[400px]">
              <div className="mx-auto w-full">
                <div className="mb-3 flex h-14 w-14 items-center justify-center rounded-2xl border border-components-panel-border-subtle font-bold shadow-lg">
                  <RiCheckboxCircleFill className='h-6 w-6 text-text-success' />
                </div>
                <h2 className="title-4xl-semi-bold text-text-primary">
                  {t('login.registerSuccess') || '注册成功'}
                </h2>
              </div>
              <div className="mx-auto mt-6 w-full">
                <Button variant='primary' className='w-full' onClick={() => {
                  setLeftTime(undefined)
                  router.replace(getSignInUrl())
                }}>{t('login.goToLogin') || '前往登录'} ({Math.round(countdown / 1000)}) </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default RegisterSetInfoPage
