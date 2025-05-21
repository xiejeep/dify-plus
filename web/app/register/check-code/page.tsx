'use client'
import Link from 'next/link'
import { RiArrowLeftLine, RiMailLine } from '@remixicon/react'
import { useTranslation } from 'react-i18next'
import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useContext } from 'use-context-selector'
import { useCountDown } from 'ahooks'
import { COUNT_DOWN_KEY, COUNT_DOWN_TIME_MS } from '../../components/signin/countdown'
import Button from '@/app/components/base/button'
import Input from '@/app/components/base/input'
import Toast from '@/app/components/base/toast'
import { sendRegisterCode, verifyRegisterCode } from '@/service/common'
import I18NContext from '@/context/i18n'

const RegisterCheckCodePage = () => {
  const { t } = useTranslation()
  const router = useRouter()
  const searchParams = useSearchParams()
  const { locale } = useContext(I18NContext)
  const email = decodeURIComponent(searchParams.get('email') || '')
  const token = decodeURIComponent(searchParams.get('token') || '')
  const [code, setCode] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [targetDate, setTargetDate] = useState<number>()
  const [countdown] = useCountDown({
    targetDate,
    onEnd: () => {
      localStorage.removeItem(COUNT_DOWN_KEY)
    },
  })

  useEffect(() => {
    const countDownTimeStamp = localStorage.getItem(COUNT_DOWN_KEY)
    if (countDownTimeStamp) {
      const targetTimeStamp = parseInt(countDownTimeStamp)
      const nowTimeStamp = Date.now()
      if (targetTimeStamp > nowTimeStamp)
        setTargetDate(targetTimeStamp)
    }
  }, [])

  const verify = async () => {
    try {
      if (!code.trim()) {
        Toast.notify({
          type: 'error',
          message: t('login.checkCode.emptyCode'),
        })
        return
      }
      if (!/\d{6}/.test(code)) {
        Toast.notify({
          type: 'error',
          message: t('login.checkCode.invalidCode'),
        })
        return
      }
      setIsLoading(true)
      const ret = await verifyRegisterCode({ email, code, token })
      ret.is_valid && router.push(`/register/set-info?${searchParams.toString()}`)
    }
    catch (error) { console.error(error) }
    finally {
      setIsLoading(false)
    }
  }

  const resendCode = async () => {
    try {
      const res = await sendRegisterCode(email, locale)
      if (res.result === 'success') {
        localStorage.setItem(COUNT_DOWN_KEY, `${COUNT_DOWN_TIME_MS}`)
        setTargetDate(Date.now() + COUNT_DOWN_TIME_MS)
        const params = new URLSearchParams(searchParams)
        params.set('token', encodeURIComponent(res.data))
        router.replace(`/register/check-code?${params.toString()}`)
      }
    }
    catch (error) { console.error(error) }
  }

  return (
    <div className='flex min-h-screen flex-col bg-white'>
      <div className='flex flex-1 flex-col items-center justify-center'>
        <div className='w-full max-w-[400px] px-4'>
          <div className='mb-8 flex items-center justify-between'>
            <Link href='/register' className='group flex items-center text-sm font-normal text-gray-500 hover:text-gray-900'>
              <RiArrowLeftLine className='mr-1 h-4 w-4 text-gray-500 group-hover:text-gray-900' />
              {t('login.back')}
            </Link>
          </div>
          <div className='mb-6 flex items-center'>
            <div className='flex h-10 w-10 items-center justify-center rounded-lg border border-gray-100 shadow-xs'>
              <RiMailLine className='h-5 w-5 text-primary-600' />
            </div>
            <div className='ml-3'>
              <div className='text-xl font-semibold text-gray-900'>{t('login.checkCode.title')}</div>
              <div className='text-sm text-gray-500'>{t('login.checkCode.description')}</div>
            </div>
          </div>
          <div className='mb-5 rounded-lg border border-gray-100 bg-gray-50 p-3'>
            <div className='text-sm text-gray-700'>{t('login.checkCode.tip')}</div>
            <div className='mt-1 text-sm font-medium text-gray-900'>{email}</div>
          </div>
          <div className='mb-4'>
            <label htmlFor="code" className='system-md-semibold my-2 text-text-secondary'>{t('login.checkCode.code')}</label>
            <div className='mt-1'>
              <Input id='code' type="text" autoComplete='off' disabled={isLoading} value={code} placeholder={t('login.checkCode.codePlaceholder') as string} onChange={e => setCode(e.target.value)} />
            </div>
          </div>
          <div className='mb-4'>
            <Button loading={isLoading} disabled={isLoading} variant='primary' className='w-full' onClick={verify}>{t('login.checkCode.verify')}</Button>
          </div>
          <div className='mb-4 flex items-center justify-center'>
            <div className='text-xs text-gray-500'>{t('login.checkCode.notReceive')}</div>
            {countdown > 0
              ? <div className='ml-1 text-xs text-gray-500'>{t('login.checkCode.resend')} ({Math.round(countdown / 1000)}s)</div>
              : <div onClick={resendCode} className='ml-1 cursor-pointer text-xs text-primary-600 hover:text-primary-700'>{t('login.checkCode.resend')}</div>
            }
          </div>
        </div>
      </div>
    </div>
  )
}

export default RegisterCheckCodePage
