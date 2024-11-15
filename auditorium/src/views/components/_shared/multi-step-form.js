/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

/** @jsx h */
const { h } = require('preact')
const { useState, useRef, useEffect } = require('preact/hooks')
const { forwardRef } = require('preact/compat')

const MultiStepForm = (props) => {
  const { onsubmit, children, steps, ...otherProps } = props

  const autofocusRef = useRef(null)
  const [step, setStep] = useState(0)
  const formData = useRef({})

  useEffect(function focusRef () {
    if (autofocusRef.current) {
      autofocusRef.current.focus()
    }
  }, [step])

  function handleCancel () {
    formData.current = {}
    setStep(0)
  }

  function handleSubmit (e) {
    e.preventDefault()
    const submission = new window.FormData(e.currentTarget)
    const values = {}
    submission.forEach((value, key) => {
      values[key] = value
    })
    Object.assign(formData.current, values)
    if (step < steps.length - 1) {
      setStep(step + 1)
    } else {
      onsubmit(formData.current, handleCancel)
    }
  }

  const current = forwardRef(steps[step])({
    onCancel: handleCancel,
    ref: autofocusRef
  })
  current.key = step

  return (
    <form
      {...otherProps}
      onsubmit={handleSubmit}
    >
      {current}
    </form>
  )
}

module.exports = MultiStepForm
