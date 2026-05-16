import { describe, expect, it } from 'vitest'

import en from '../locales/en'
import ptBR from '../locales/pt-BR'
import zh from '../locales/zh'

describe('usage service tier locale keys', () => {
  it('contains zh labels for service tier tooltip', () => {
    expect(zh.usage.serviceTier).toBe('服务档位')
    expect(zh.usage.serviceTierPriority).toBe('Fast')
    expect(zh.usage.serviceTierFlex).toBe('Flex')
    expect(zh.usage.serviceTierStandard).toBe('Standard')
  })

  it('contains en labels for service tier tooltip', () => {
    expect(en.usage.serviceTier).toBe('Service tier')
    expect(en.usage.serviceTierPriority).toBe('Fast')
    expect(en.usage.serviceTierFlex).toBe('Flex')
    expect(en.usage.serviceTierStandard).toBe('Standard')
  })

  it('contains pt-BR labels and falls back to inherited English keys', () => {
    expect(ptBR.usage.serviceTier).toBe('Classe de serviço')
    expect(ptBR.usage.serviceTierPriority).toBe('Rápido')
    expect(ptBR.usage.serviceTierFlex).toBe('Flex')
    expect(ptBR.usage.serviceTierStandard).toBe('Padrão')
    expect(ptBR.keyUsage.title).toBe(en.keyUsage.title)
  })
})
