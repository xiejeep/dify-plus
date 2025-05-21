'use client'
import Link from 'next/link'
import { RiArrowLeftLine, RiUserAddLine } from '@remixicon/react'
import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useContext } from 'use-context-selector'
import { COUNT_DOWN_KEY, COUNT_DOWN_TIME_MS } from '../components/signin/countdown'
import { emailRegex } from '@/config'
import Button from '@/app/components/base/button'
import Input from '@/app/components/base/input'
import Toast from '@/app/components/base/toast'
import { sendRegisterCode } from '@/service/common'
import I18NContext from '@/context/i18n'
import { noop } from 'lodash-es'

const RegisterPage = () => {
  const { t } = useTranslation()
  const router = useRouter()
  const searchParams = useSearchParams()
  const { locale } = useContext(I18NContext)
  const [email, setEmail] = useState(decodeURIComponent(searchParams.get('email') || ''))
  const [loading, setLoading] = useState(false)

  const handleGetEMailVerificationCode = async () => {
    try {
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
      setLoading(true)
      const res = await sendRegisterCode(email, locale)
      if (res.result === 'success') {
        localStorage.setItem(COUNT_DOWN_KEY, `${COUNT_DOWN_TIME_MS}`)
        const params = new URLSearchParams(searchParams)
        params.set('token', encodeURIComponent(res.data))
        params.set('email', encodeURIComponent(email))
        router.push(`/register/check-code?${params.toString()}`)
      } else if (res.code === 'account_already_exists') {
        Toast.notify({
          type: 'error',
          message: t('login.error.accountAlreadyExists') || '账号已存在，请直接登录',
        })
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className='flex min-h-screen flex-col bg-white'>
      <div className='flex flex-1 flex-col items-center justify-center'>
        <div className='w-full max-w-[400px] px-4'>
          <div className='mb-8 flex items-center justify-between'>
            <Link href='/signin' className='group flex items-center text-sm font-normal text-gray-500 hover:text-gray-900'>
              <RiArrowLeftLine className='mr-1 h-4 w-4 text-gray-500 group-hover:text-gray-900' />
              {t('login.back')}
            </Link>
          </div>
          <div className='mb-6 flex items-center'>
            <div className='flex h-10 w-10 items-center justify-center rounded-lg border border-gray-100 shadow-xs'>
              <RiUserAddLine className='h-5 w-5 text-primary-600' />
            </div>
            <div className='ml-3'>
              <div className='text-xl font-semibold text-gray-900'>{t('login.register') || '注册'}</div>
              <div className='text-sm text-gray-500'>{t('login.registerTip') || '创建您的账号'}</div>
            </div>
          </div>
          <form onSubmit={noop}>
            <input type='text' className='hidden' />
            <div className='mb-2'>
              <label htmlFor="email" className='system-md-semibold my-2 text-text-secondary'>{t('login.email')}</label>
              <div className='mt-1'>
                <Input id='email' type="email" disabled={loading} value={email} placeholder={t('login.emailPlaceholder') as string} onChange={e => setEmail(e.target.value)} />
              </div>
              <div className='mt-3'>
                <Button loading={loading} disabled={loading} variant='primary' className='w-full' onClick={handleGetEMailVerificationCode}>{t('login.sendVerificationCode')}</Button>
              </div>
            </div>
          </form>
          <div className='py-2'>
            <div className='h-px bg-gradient-to-r from-background-gradient-mask-transparent via-divider-regular to-background-gradient-mask-transparent'></div>
          </div>
          <div className='text-center text-xs text-gray-500'>
            {t('login.registerAgreement') || '注册即表示您同意我们的'} <Link href="/terms" className='text-primary-600 hover:text-primary-700'>{t('login.terms') || '服务条款'}</Link> {t('login.and') || '和'} <Link href="/privacy" className='text-primary-600 hover:text-primary-700'>{t('login.privacy') || '隐私政策'}</Link>
          </div>
        </div>
      </div>
    </div>
  )
}

export default RegisterPage
