<template>
  <!-- Admin custom content mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- NexusMind Landing Page -->
  <div v-else class="nm-home">
    <div class="nm-bg-grid"></div>
    <div class="nm-bg-glow"></div>

    <!-- ── Header ─────────────────────────────────────────── -->
    <header class="nm-header">
      <div class="nm-wrap nm-header-inner">
        <a href="/home" class="nm-brand">
          <div class="nm-brand-logo">
            <img :src="siteLogo || '/logo.png'" alt="NexusMind" />
          </div>
          <div class="nm-brand-name">
            <b>NexusMind</b>
            <span>AI GATEWAY · BR</span>
          </div>
        </a>

        <nav class="nm-nav">
          <a href="#nm-how" class="nm-nav-link">$ como_funciona</a>
          <a href="#nm-math" class="nm-nav-link">$ matemática</a>
          <a href="#nm-pricing" class="nm-nav-link">$ preço</a>
          <a href="#nm-faq" class="nm-nav-link">$ faq</a>
        </nav>

        <div class="nm-actions">
          <LocaleSwitcher />
          <router-link v-if="isAuthenticated" :to="dashboardPath" class="nm-btn nm-btn-ghost">
            {{ t('home.dashboard') }}
          </router-link>
          <router-link v-else to="/login" class="nm-btn nm-btn-ghost">
            {{ t('home.login') }}
          </router-link>
          <a href="#nm-pricing" class="nm-btn nm-btn-primary">
            Começar <span class="nm-arr">→</span>
          </a>
        </div>
      </div>
    </header>

    <!-- ── Main ───────────────────────────────────────────── -->
    <main>

      <!-- Hero -->
      <section id="nm-top" class="nm-hero nm-wrap">
        <div class="nm-hero-grid">
          <div>
            <div class="nm-eyebrow">
              <span class="nm-eyebrow-dot"></span>
              <span>NexusMind · <span class="nm-eyebrow-v">v2.1</span> · BR</span>
            </div>
            <h1 class="nm-h1">
              Uma chave.<br />
              <span class="nm-accent">Todos os modelos</span><br />
              de IA.
            </h1>
            <p class="nm-lede">
              Roteamento inteligente para Claude, GPT e Gemini —
              no real, com Pix. Painel pra monitorar tokens, custo
              e failover automático quando um provider engasga.
            </p>
            <div class="nm-cta-row">
              <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="nm-btn nm-btn-primary nm-btn-lg">
                {{ isAuthenticated ? 'Abrir painel' : 'Quero minha chave' }}
                <span class="nm-arr">→</span>
              </router-link>
              <a href="#nm-how" class="nm-btn nm-btn-ghost nm-btn-lg">Como funciona</a>
            </div>
            <div class="nm-pill-row">
              <span class="nm-pill-item">
                <span class="nm-dot-green"></span> 100+ modelos · múltiplos providers
              </span>
              <span class="nm-pill-item">
                <span class="nm-ok">●</span> 99.94% uptime · 30d
              </span>
              <span class="nm-pill-item">★ 4.9 · 100+ reviews</span>
            </div>
          </div>

          <!-- Routing console -->
          <div class="nm-console">
            <div class="nm-console-head">
              <span class="nm-console-dots"><i /><i /><i /></span>
              <span class="nm-console-title">gateway.nexusmind.digital</span>
              <span class="nm-console-live">
                <i class="nm-live-dot"></i> ROUTING
              </span>
            </div>
            <div class="nm-console-body">
              <!-- Routing SVG diagram -->
              <svg class="nm-diagram" viewBox="0 0 560 320" preserveAspectRatio="xMidYMid meet">
                <defs>
                  <linearGradient id="nmG1" x1="0" y1="0" x2="1" y2="0">
                    <stop offset="0%" stop-color="#5fd5ff" />
                    <stop offset="100%" stop-color="#a78bfa" />
                  </linearGradient>
                  <radialGradient id="nmCoreGlow" cx="50%" cy="50%" r="50%">
                    <stop offset="0%" stop-color="rgba(95,213,255,0.45)" />
                    <stop offset="100%" stop-color="rgba(95,213,255,0)" />
                  </radialGradient>
                </defs>
                <g opacity="0.07">
                  <line v-for="n in 7" :key="n" x1="0" :y1="n*46+20" x2="560" :y2="n*46+20" stroke="#5fd5ff" stroke-dasharray="2 6" />
                </g>
                <text x="14" y="22" class="nm-svg-dim">$ ingress</text>
                <text x="540" y="22" class="nm-svg-dim" text-anchor="end">$ providers</text>

                <!-- Request nodes -->
                <g v-for="(req, i) in nmRequests" :key="i" :opacity="activeRoute === i ? 1 : 0.55">
                  <rect :x="14" :y="req.y-12" width="180" height="24" rx="4" fill="#0a0d15" :stroke="activeRoute === i ? '#5fd5ff' : 'rgba(255,255,255,0.11)'" />
                  <circle cx="26" :cy="req.y" r="3" :fill="activeRoute === i ? '#5fd5ff' : '#a78bfa'" />
                  <text x="36" :y="req.y+4" :class="['nm-svg-label', activeRoute === i ? 'nm-svg-cyan' : 'nm-svg-dim']">{{ req.label }}</text>
                </g>

                <!-- Lanes to core -->
                <path
                  v-for="(req, i) in nmRequests"
                  :key="'l'+i"
                  :d="`M194 ${req.y} C 240 ${req.y}, 240 160, 280 160`"
                  stroke="url(#nmG1)"
                  :class="['nm-lane', activeRoute === i ? 'nm-lane-active' : '']"
                />

                <!-- Core (router) -->
                <circle cx="280" cy="160" r="58" fill="url(#nmCoreGlow)" />
                <rect x="240" y="138" width="80" height="44" rx="6" fill="#0d1119" stroke="#5fd5ff" />
                <text x="280" y="155" class="nm-svg-cyan" text-anchor="middle" font-size="11" font-weight="600" font-family="JetBrains Mono,monospace" style="letter-spacing: 0.06em">NEXUS</text>
                <text x="280" y="170" class="nm-svg-dim" text-anchor="middle">nexusmind.digital</text>

                <!-- Provider lanes + boxes -->
                <g v-for="(prov, i) in nmProviders" :key="'p'+i">
                  <path
                    :d="`M320 160 C 380 160, 380 ${prov.y}, 420 ${prov.y}`"
                    stroke="url(#nmG1)"
                    :class="['nm-lane', activeRoute === i ? 'nm-lane-active' : '']"
                  />
                  <rect :x="420" :y="prov.y-12" width="126" height="24" rx="4" :class="['nm-prov-svg', activeRoute === i ? 'nm-prov-svg-active' : '']" />
                  <rect :x="426" :y="prov.y-7" width="14" height="14" rx="2" :fill="activeRoute === i ? '#5fd5ff' : '#1c2233'" />
                  <text :x="448" :y="prov.y+4" :class="['nm-svg-label', activeRoute === i ? 'nm-svg-cyan' : 'nm-svg-dim']">{{ prov.name }}</text>
                  <text x="540" :y="prov.y+4" class="nm-svg-dim" text-anchor="end" font-size="9">{{ activeRoute === i ? '200 OK' : '—' }}</text>
                </g>

                <circle class="nm-packet nm-packet-1" r="4" />
                <circle class="nm-packet nm-packet-2" r="4" />
                <circle class="nm-packet nm-packet-3" r="4" />
              </svg>

              <!-- Routing target tag -->
              <div class="nm-routing-target">
                <span class="nm-rt-label">$ now routing</span>
                <span class="nm-rt-val">{{ nmRequests[activeRoute].method }}</span>
                <span class="nm-rt-label">→</span>
                <span class="nm-rt-val">{{ nmRequests[activeRoute].prov }}</span>
                <span class="nm-rt-val">{{ nmRequests[activeRoute].model }}</span>
              </div>

              <!-- Live log feed -->
              <div class="nm-log-feed">
                <div v-for="(log, i) in logs" :key="log.id" class="nm-log-line" :style="{ opacity: 1 - i * 0.32 }">
                  <span class="nm-log-t">[{{ log.t }}]</span>
                  <span class="nm-log-method">{{ log.path }}</span>
                  <span class="nm-log-arr">→</span>
                  <span class="nm-log-prov">{{ log.prov }}</span>
                  <span class="nm-log-model">{{ log.model }}</span>
                  <span class="nm-log-latency">{{ log.latency }}ms</span>
                  <span class="nm-log-ok">200</span>
                </div>
              </div>

              <div class="nm-console-meta">
                <div class="nm-meta-cell">
                  <div class="nm-meta-k">catálogo</div>
                  <div class="nm-meta-v">100+<span class="nm-meta-u">modelos</span></div>
                </div>
                <div class="nm-meta-cell">
                  <div class="nm-meta-k">latência p95</div>
                  <div class="nm-meta-v">312<span class="nm-meta-u">ms</span></div>
                </div>
                <div class="nm-meta-cell">
                  <div class="nm-meta-k">economia média</div>
                  <div class="nm-meta-v">R$ 851<span class="nm-meta-u">/mês</span></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ── Tutorial ─────────────────────────────────────── -->
      <section id="nm-tutorial" class="nm-section nm-wrap">
        <div class="nm-s-label nm-reveal"><span class="nm-sq"></span> $ tutorial <span class="nm-tic">--duration=90s</span></div>
        <h2 class="nm-s-title nm-reveal">Veja em 90 segundos.</h2>
        <p class="nm-s-sub nm-reveal">Do Pix até a primeira chamada. Sem áudio, sem firula — só os cliques na tela.</p>

        <div class="nm-tut-grid">
          <div class="nm-reveal">
            <div class="nm-video-wrap">
              <iframe
                v-if="videoEmbedUrl"
                :src="videoEmbedUrl"
                title="Tutorial NexusMind"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowfullscreen
              ></iframe>
              <div v-else class="nm-video-placeholder">
                <div class="nm-video-badge">✦ TUTORIAL · 90s</div>
                <div class="nm-video-time">SEM ÁUDIO</div>
                <button class="nm-play-btn" aria-label="Reproduzir tutorial">
                  <svg viewBox="0 0 24 24" width="32" height="32" fill="currentColor">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </button>
                <div class="nm-video-hint">Cole o link do vídeo em <code>videoEmbedUrl</code></div>
                <div class="nm-video-strip">
                  <span class="nm-vs-left"><span class="nm-vs-dot"></span> NEXUSMIND · /tutorial</span>
                  <span class="nm-vs-progress"></span>
                  <span>00:22 / 01:30</span>
                </div>
              </div>
            </div>
          </div>

          <div class="nm-tut-chapters nm-reveal">
            <h4 class="nm-tut-h">$ capítulos --chapters</h4>
            <ul class="nm-tut-list">
              <li v-for="(c, i) in nmChapters" :key="i">
                <span class="nm-tut-time">{{ c.t }}</span>
                <span>{{ c.txt }}</span>
              </li>
            </ul>
          </div>
        </div>
      </section>

      <!-- Como funciona -->
      <section id="nm-how" class="nm-section nm-wrap">
        <div class="nm-s-label"><span class="nm-sq"></span> $ como_funciona <span class="nm-tic">--steps=3</span></div>
        <h2 class="nm-s-title">Em 3 passos.<br />Sem firula.</h2>
        <p class="nm-s-sub">Você não precisa migrar nada, criar conta em provider gringo, nem pagar dólar. Só trocar a base URL.</p>

        <div class="nm-steps">
          <article class="nm-step">
            <div class="nm-step-num"><b>01</b> /// CRIAR CHAVE</div>
            <h3>Pagar e receber a chave</h3>
            <p>Pix confirma na hora. A chave aparece no seu painel, no mesmo email da compra.</p>
            <div class="nm-codeblock">
              <span class="nm-c"># painel.nexusmind.digital</span><br />
              <span class="nm-k">sk_live</span>_<span class="nm-v">7c9f...8b2a</span>
            </div>
          </article>
          <article class="nm-step">
            <div class="nm-step-num"><b>02</b> /// CONFIGURAR</div>
            <h3>Trocar a base URL</h3>
            <p>Cole a chave e a base URL no seu Cursor, Claude Code, script ou SDK. Funciona em qualquer client OpenAI-compat.</p>
            <div class="nm-codeblock">
              <span class="nm-k">ANTHROPIC_BASE_URL</span>=<br />
              <span class="nm-s">"https://gateway.nexusmind.digital"</span>
            </div>
          </article>
          <article class="nm-step">
            <div class="nm-step-num"><b>03</b> /// USAR</div>
            <h3>Pedir e roteamos</h3>
            <p>Seu request chega no gateway, escolhemos o provider saudável mais barato pro modelo pedido e devolvemos a resposta.</p>
            <div class="nm-codeblock">
              <span class="nm-c">→ model=claude-opus-4-7</span><br />
              <span class="nm-c">→ routed to provider AN_03</span><br />
              <span class="nm-s">← 200 OK · 482ms</span>
            </div>
          </article>
        </div>
      </section>

      <!-- Matemática -->
      <section id="nm-math" class="nm-section nm-wrap">
        <div class="nm-s-label"><span class="nm-sq"></span> $ matemática <span class="nm-tic">--why=cheaper</span></div>
        <h2 class="nm-s-title">Por que custa <span class="nm-accent">R$ 249</span><br />no lugar de R$ 1.020?</h2>
        <p class="nm-s-sub">Sem pegadinha. Mesmos modelos, mesma API, mesma performance. A diferença vem de três fatores.</p>

        <div class="nm-math-grid">
          <div class="nm-compare">
            <div class="nm-cmp-row nm-cmp-header">
              <div class="nm-cmp-cell">Assinatura mensal</div>
              <div class="nm-cmp-cell">Direto no provider</div>
              <div class="nm-cmp-cell nm-cmp-ours">Via NexusMind</div>
            </div>
            <div v-for="row in nmMathRows" :key="row.label" class="nm-cmp-row">
              <div class="nm-cmp-cell nm-cmp-label">{{ row.label }}</div>
              <div class="nm-cmp-cell"><span class="nm-strike">{{ row.direct }}</span></div>
              <div class="nm-cmp-cell nm-cmp-ours">incluso</div>
            </div>
            <div class="nm-cmp-row nm-cmp-total">
              <div class="nm-cmp-cell">Total / mês</div>
              <div class="nm-cmp-cell"><span class="nm-strike">R$ 1.020</span></div>
              <div class="nm-cmp-cell nm-cmp-ours">R$ 249</div>
            </div>
          </div>

          <div class="nm-math-points">
            <div
              v-for="(pt, i) in nmMathPoints"
              :key="i"
              class="nm-math-pt"
              :class="{ 'nm-math-pt-active': activeMath === i }"
              @mouseenter="activeMath = i"
            >
              <div class="nm-pt-tag">{{ pt.tag }}</div>
              <h4>{{ pt.title }}</h4>
              <p>{{ pt.body }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ── Integrações ──────────────────────────────────── -->
      <section id="nm-integrations" class="nm-section nm-wrap">
        <div class="nm-s-label nm-reveal"><span class="nm-sq"></span> $ integrações <span class="nm-tic">--openai-compat=true</span></div>
        <h2 class="nm-s-title nm-reveal">Cola a key onde<br />você já trabalha.</h2>
        <p class="nm-s-sub nm-reveal">Trocou a base URL e a chave, funcionou. Qualquer client OpenAI-compat funciona — sem migração, sem refactor.</p>

        <div class="nm-int-grid">
          <article
            v-for="(t, i) in nmIntegrations"
            :key="i"
            class="nm-int-card nm-reveal"
            :style="{ transitionDelay: (i * 40) + 'ms' }"
          >
            <span v-if="t.tested" class="nm-int-tested"><span class="nm-int-d"></span> Testado</span>
            <div class="nm-int-row">
              <div class="nm-int-icon">{{ t.icon }}</div>
              <div class="nm-int-name">{{ t.name }}</div>
            </div>
            <div class="nm-int-where">{{ t.where }}</div>
            <div class="nm-int-env"><b>{{ t.env }}</b>=<span class="nm-int-url">gateway.nexusmind.digital</span></div>
          </article>
        </div>

        <div class="nm-int-more nm-reveal">
          <b>+ também roda em:</b> Zed · Vercel AI SDK · LlamaIndex · OpenRouter clients · Goose · Roo Code · Cody · scripts Python/Node — qualquer ferramenta que aceite uma chave compatível com OpenAI.
        </div>
      </section>

      <!-- Providers -->
      <section id="nm-providers" class="nm-section nm-wrap">
        <div class="nm-s-label"><span class="nm-sq"></span> $ providers <span class="nm-tic">--list --all</span></div>
        <h2 class="nm-s-title">Mais de 100 modelos<br />na mesma chave.</h2>
        <p class="nm-s-sub">Catálogo amplo para coding, reasoning, multimodal e open-source. Abaixo ficam as famílias principais — novos modelos entram no roteador sem você trocar integração.</p>
        <div class="nm-providers">
          <article v-for="pv in nmProviderCards" :key="pv.code" class="nm-prov-card">
            <span class="nm-prov-status"><span class="nm-prov-dot"></span> ON</span>
            <div class="nm-prov-ico" :style="{ color: pv.color, borderColor: pv.color + '33' }">{{ pv.code }}</div>
            <h4>{{ pv.name }}</h4>
            <div class="nm-prov-models">
              <b>{{ pv.label }}</b>
              <span v-for="model in pv.models" :key="model">{{ model }}</span>
            </div>
          </article>
        </div>
        <div class="nm-model-note">
          <b>Catálogo dinâmico:</b> além dos modelos acima, o gateway expõe combos como balanced-load, high-availability, coding, reasoning, economy, opus, flash e raiko.
        </div>
      </section>

      <!-- Pricing -->
      <section id="nm-pricing" class="nm-section nm-wrap">
        <div class="nm-s-label"><span class="nm-sq"></span> $ pricing <span class="nm-tic">--currency=BRL</span></div>
        <h2 class="nm-s-title">Três planos.<br />Zero letra miúda.</h2>
        <p class="nm-s-sub">Cancela quando quiser. Pagamento facilitado via Pix ou cartão.</p>
        <div class="nm-pricing-grid">
          <article
            v-for="plan in nmPlans"
            :key="plan.name"
            class="nm-plan"
            :class="{ 'nm-plan-feat': plan.feat }"
          >
            <span v-if="plan.feat" class="nm-plan-badge">Recomendado</span>
            <div class="nm-plan-type">{{ plan.type }}</div>
            <h3>{{ plan.name }}</h3>
            <p class="nm-plan-desc">{{ plan.desc }}</p>
            <div class="nm-plan-price">
              <span class="nm-price-cur">R$</span>
              <span class="nm-price-amt">{{ plan.price }}</span>
              <span class="nm-price-per">/mês</span>
            </div>
            <div v-if="plan.oldPrice" class="nm-plan-old">Equivalente direto · R$ {{ plan.oldPrice }}/mês</div>
            <div v-else class="nm-plan-old nm-plan-old-spacer">·</div>
            <ul class="nm-plan-features">
              <li v-for="f in plan.features" :key="f">{{ f }}</li>
              <li v-for="f in plan.muted" :key="f" class="nm-feat-muted">{{ f }}</li>
            </ul>
            <router-link
              to="/purchase?tab=subscription"
              class="nm-btn nm-btn-lg"
              :class="plan.feat ? 'nm-btn-primary' : 'nm-btn-dark'"
              style="justify-content: center; display: flex;"
            >
              {{ plan.feat ? 'Assinar Pro' : 'Escolher ' + plan.name }}
              <span class="nm-arr">→</span>
            </router-link>
          </article>
        </div>
      </section>

      <!-- Proof -->
      <section id="nm-proof" class="nm-section nm-wrap">
        <div class="nm-s-label"><span class="nm-sq"></span> $ devs --aprovam <span class="nm-tic">--count=100+</span></div>
        <h2 class="nm-s-title">Devs brasileiros<br />respirando IA o dia inteiro.</h2>
        <p class="nm-s-sub">Quem usa NexusMind tá fazendo refactor pesado, agentes autônomos e automações em loop. Não dá pra parar.</p>
        <div class="nm-proof-grid">
          <article v-for="(tm, i) in nmTestimonials" :key="i" class="nm-tcard">
            <div class="nm-tcard-meta">
              <span>{{ tm.meta[0] }}</span>
              <span class="nm-t-sep">/</span>
              <span class="nm-t-ec">{{ tm.meta[1] }}</span>
            </div>
            <p class="nm-tcard-quote">"{{ tm.quote }}"</p>
            <div class="nm-tcard-who">
              <div class="nm-tcard-av">{{ tm.initial }}</div>
              <div>
                <b>{{ tm.name }}</b>
                <span>{{ tm.role }}</span>
              </div>
            </div>
          </article>
        </div>
      </section>

      <!-- FAQ -->
      <section id="nm-faq" class="nm-section nm-wrap">
        <div class="nm-s-label"><span class="nm-sq"></span> $ faq <span class="nm-tic">--top=8</span></div>
        <h2 class="nm-s-title">Dúvidas que<br />aparecem no DM.</h2>
        <p class="nm-s-sub">Se a sua não tá aqui, chama no WhatsApp que respondo direto. Sem bot.</p>
        <div class="nm-faq-grid">
          <div
            v-for="(faq, i) in nmFaqs"
            :key="i"
            class="nm-faq-item"
            :class="{ 'nm-faq-open': openFaq === i }"
          >
            <button @click="toggleFaq(i)">
              <span>
                <span class="nm-faq-num">Q.{{ String(i + 1).padStart(2, '0') }}</span>
                {{ faq.q }}
              </span>
              <span class="nm-faq-plus">+</span>
            </button>
            <div class="nm-faq-answer"><p>{{ faq.a }}</p></div>
          </div>
        </div>
      </section>

      <!-- Final CTA -->
      <section class="nm-section nm-wrap">
        <div class="nm-final-cta">
          <div class="nm-s-label" style="justify-content: center;">
            <span class="nm-sq"></span> $ próximo_passo <span class="nm-tic">--ready</span>
          </div>
          <h2 class="nm-s-title">Pronto pra rotear seus tokens?</h2>
          <p class="nm-cta-p">Pix confirma em segundos. Em 4 minutos você tá rodando Claude Opus 4.7 no seu Cursor, pagando em real.</p>
          <div class="nm-cta-btns">
            <router-link to="/purchase?tab=subscription" class="nm-btn nm-btn-primary nm-btn-lg">
              Quero minha chave <span class="nm-arr">→</span>
            </router-link>
            <a :href="contactUrl" class="nm-btn nm-btn-ghost nm-btn-lg">Falar no WhatsApp</a>
          </div>
        </div>
      </section>

    </main>

    <!-- ── Footer ──────────────────────────────────────────── -->
    <footer class="nm-footer">
      <div class="nm-wrap nm-footer-grid">
        <div class="nm-footer-col nm-footer-about">
          <div class="nm-brand">
            <div class="nm-brand-logo">
              <img :src="siteLogo || '/logo.png'" alt="" />
            </div>
            <div class="nm-brand-name">
              <b>{{ siteName }}</b>
              <span>AI GATEWAY · BR</span>
            </div>
          </div>
          <p class="nm-footer-desc">
            Gateway de roteamento e cota pra APIs de IA. Acesse Claude, GPT e Gemini
            com uma chave só, em real, no Pix. Não somos a Anthropic — somos seu intermediário brasileiro.
          </p>
        </div>
        <div class="nm-footer-col">
          <h5>Produto</h5>
          <a href="#nm-how">Como funciona</a>
          <a href="#nm-providers">Providers</a>
          <a href="#nm-pricing">Pricing</a>
          <router-link :to="dashboardPath">Painel</router-link>
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a>
        </div>
        <div class="nm-footer-col">
          <h5>Canais</h5>
          <a :href="contactUrl" target="_blank" rel="noopener noreferrer">WhatsApp</a>
          <a href="https://youtube.com" target="_blank" rel="noopener noreferrer">YouTube</a>
          <a href="https://instagram.com" target="_blank" rel="noopener noreferrer">Instagram</a>
        </div>
        <div class="nm-footer-col">
          <h5>Legal</h5>
          <a href="#">Termos de uso</a>
          <a href="#">Privacidade</a>
          <a href="#">Status</a>
          <a href="#">Reembolso</a>
        </div>
      </div>
      <div class="nm-wrap nm-footer-bot">
        <span>© {{ currentYear }} {{ siteName }} · Pagamento facilitado via Pix ou cartão</span>
        <span>nexusmind.digital</span>
      </div>
    </footer>

    <!-- ── Floating sticky CTA pill ── -->
    <div v-if="!floatDismissed" class="nm-float-cta" :class="{ 'nm-float-visible': floatVisible }">
      <div class="nm-float-inner">
        <span class="nm-float-prompt">{{ '>_' }} NexusMind</span>
        <span class="nm-float-sep">·</span>
        <a :href="contactUrl" target="_blank" rel="noopener noreferrer" class="nm-float-link">WhatsApp</a>
        <a href="#" class="nm-float-link">YouTube</a>
        <a href="#" class="nm-float-link">Instagram</a>
        <router-link to="/purchase?tab=subscription" class="nm-float-btn">Assinar · R$ 249 <span>→</span></router-link>
        <button class="nm-float-close" @click="floatDismissed = true" aria-label="Dispensar">×</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { useAuthStore, useAppStore } from '@/stores'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

// ── Auth & Settings ────────────────────────────────────────
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const currentYear = computed(() => new Date().getFullYear())

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'NexusMind')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const contactUrl = computed(() => appStore.cachedPublicSettings?.contact_info || 'https://wa.me/')

const isHomeContentUrl = computed(() => {
  const c = homeContent.value.trim()
  return c.startsWith('http://') || c.startsWith('https://')
})

// ── Routing diagram animation ──────────────────────────────
const activeRoute = ref(0)
interface LogEntry { id: number; t: string; path: string; prov: string; model: string; latency: number }
const logs = ref<LogEntry[]>([])

const nmRequests = [
  { y: 80, label: 'POST /v1/messages', method: 'POST /v1/messages', path: '/v1/messages', model: 'claude-opus-4.7', prov: 'AN_03' },
  { y: 140, label: 'POST /v1/chat/completions', method: 'POST /v1/chat/completions', path: '/v1/chat/completions', model: 'gpt-5', prov: 'OA_01' },
  { y: 200, label: 'POST /v1/embeddings', method: 'POST /v1/embeddings', path: '/v1/embeddings', model: 'gemini-2.5-flash', prov: 'GG_02' },
]
const nmProviders = [
  { y: 70, name: 'ANTHROPIC' },
  { y: 160, name: 'OPENAI' },
  { y: 250, name: 'GOOGLE' },
]

function makeRouteLog(routeIndex: number, id: number, secondsOffset: number): LogEntry {
  const now = new Date(Date.now() - secondsOffset * 1000)
  const t = `${String(now.getHours()).padStart(2,'0')}:${String(now.getMinutes()).padStart(2,'0')}:${String(now.getSeconds()).padStart(2,'0')}`
  const req = nmRequests[routeIndex]
  return { id, t, path: req.path, prov: req.prov, model: req.model, latency: 320 + routeIndex * 47 }
}

logs.value = [
  makeRouteLog(0, 1, 0),
  makeRouteLog(1, 2, 2),
  makeRouteLog(2, 3, 4),
]

// Floating CTA pill state
const floatVisible = ref(false)
const floatDismissed = ref(false)
function onScroll() {
  if (floatDismissed.value) return
  floatVisible.value = window.scrollY > 700
}

const videoEmbedUrl = ''

const nmChapters = [
  { t: '00:00', txt: 'Pagar via Pix em 3 cliques' },
  { t: '00:15', txt: 'Receber a chave no painel' },
  { t: '00:30', txt: 'Configurar no Cursor' },
  { t: '00:55', txt: 'Configurar no Claude Code' },
  { t: '01:15', txt: 'Primeira chamada · Opus 4.7' },
  { t: '01:25', txt: 'Monitorar uso e custo' },
]

const nmIntegrations = [
  { name: 'Cursor', icon: '▸', where: 'Settings → Models → Add base URL', env: 'OPENAI_BASE_URL', tested: true },
  { name: 'Claude Code', icon: '⌘', where: 'Variável de ambiente shell', env: 'ANTHROPIC_BASE_URL', tested: true },
  { name: 'Windsurf', icon: '≋', where: 'Settings → AI → Custom endpoint', env: 'BASE_URL', tested: true },
  { name: 'Cline', icon: '/', where: 'Extensão VS Code → Settings', env: 'OPENAI_BASE_URL', tested: true },
  { name: 'Continue.dev', icon: '→', where: '~/.continue/config.json', env: 'apiBase', tested: true },
  { name: 'Aider', icon: ':', where: '--openai-api-base flag', env: 'OPENAI_API_BASE', tested: true },
  { name: 'n8n', icon: '⬡', where: 'AI / OpenAI node → custom URL', env: 'baseURL', tested: true },
  { name: 'LangChain', icon: 'λ', where: 'ChatOpenAI(base_url=…)', env: 'base_url', tested: true },
]

// Scroll reveal observer
let revealObserver: IntersectionObserver | null = null
function setupReveal() {
  if (typeof IntersectionObserver === 'undefined') return
  revealObserver = new IntersectionObserver((entries) => {
    entries.forEach(e => {
      if (e.isIntersecting) {
        e.target.classList.add('nm-reveal-in')
        revealObserver?.unobserve(e.target)
      }
    })
  }, { threshold: 0.12, rootMargin: '0px 0px -60px 0px' })
  document.querySelectorAll('.nm-reveal').forEach(el => revealObserver?.observe(el))
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
  setTimeout(setupReveal, 30)
})
onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  if (revealObserver) revealObserver.disconnect()
})

// ── Math section ───────────────────────────────────────────
const activeMath = ref(0)
const nmMathRows = [
  { label: 'Claude Pro / Max', direct: 'R$ 580' },
  { label: 'ChatGPT Plus / Pro', direct: 'R$ 320' },
  { label: 'Gemini Advanced', direct: 'R$ 120' },
]
const nmMathPoints = [
  { tag: '01 · Fator', title: 'Compra em volume', body: 'Compramos pacotes de tokens dos providers em escala. Volume continuado virou desconto continuado — e a gente repassa pra você.' },
  { tag: '02 · Fator', title: 'Roteamento inteligente', body: 'Pra cada request, escolhemos o provider com a melhor latência e o menor custo no momento. Você paga sempre o melhor preço, sem decidir nada.' },
  { tag: '03 · Fator', title: 'Sem múltiplas assinaturas', body: 'Em vez de assinar Claude Max + ChatGPT Pro + Gemini Ultra (R$ 1.020+), você paga um plano só e usa todos. A mesma performance, em real.' },
]

// ── Providers ─────────────────────────────────────────────
const nmProviderCards = [
  { code: 'CL', name: 'Claude / Anthropic', label: 'Família Claude 4.x', models: ['Opus 4.7', 'Opus 4.6', 'Sonnet 4.6', 'Sonnet 4.5', 'Haiku 4.5', 'Thinking'], color: '#d97757' },
  { code: 'CX', name: 'GPT / Codex', label: 'Família GPT 5.x + coding', models: ['GPT-5.5', 'GPT-5.5 xHigh', 'GPT-5.5 High', 'GPT-5.4', 'GPT-5.4 Mini', 'Codex Spark', 'Auto Review'], color: '#10a37f' },
  { code: 'GM', name: 'Gemini / Google', label: 'Pro, Flash, Lite e multimodal', models: ['Gemini 3.1 Pro', 'Gemini 3 Pro', 'Gemini 3 Flash', 'Gemini 2.5 Pro', '2.5 Flash', 'Flash Lite', 'Image'], color: '#4285f4' },
  { code: 'KR', name: 'Kiro / Agentic', label: 'Modelos para agentes e IDE', models: ['Auto Kiro', 'Claude Opus 4.7', 'Claude Sonnet 4.6', 'DeepSeek 3.2', 'MiniMax M2.5', 'GLM-5', 'Qwen3 Coder'], color: '#8b5cf6' },
  { code: 'OS', name: 'Open-source / Reasoning', label: 'DeepSeek, Qwen, Llama, Kimi, GLM', models: ['DeepSeek V3.2', 'DeepSeek R1', 'Qwen3 235B', 'Qwen3 Coder', 'Llama 4 Scout', 'Kimi K2.5', 'GLM-4.7'], color: '#f97316' },
  { code: 'MX', name: 'Mistral / Groq / Extras', label: 'Velocidade, áudio e modelos europeus', models: ['Mistral Large', 'Mistral Medium', 'Codestral', 'Devstral', 'Groq Llama', 'GPT-OSS 120B', 'Whisper'], color: '#ec4899' },
]

// ── Plans ─────────────────────────────────────────────────
const nmPlans = [
  {
    type: 'Starter', name: 'Solo dev', feat: false, price: 79, oldPrice: null,
    desc: 'Pra quem tá começando ou usando IA pra projetos pessoais.',
    features: ['1 chave de API', 'Acesso a Claude, GPT, Gemini', 'Limite mensal de 50M tokens', 'Painel de uso em tempo real', 'Suporte por email'],
    muted: ['Roteamento prioritário', 'Multi-time'],
  },
  {
    type: 'Recomendado', name: 'Pro', feat: true, price: 249, oldPrice: 1020,
    desc: 'Acesso completo, todos os modelos premium, sem letra miúda.',
    features: ['Chaves ilimitadas por projeto', 'Todos os modelos premium liberados', 'Limite mensal de 500M tokens', 'Roteamento prioritário entre providers', 'Failover automático entre contas', 'Logs e relatórios exportáveis', 'Suporte direto via WhatsApp'],
    muted: [],
  },
  {
    type: 'Time', name: 'Equipe', feat: false, price: 749, oldPrice: null,
    desc: 'Pra times de eng e operações de IA em escala.',
    features: ['Tudo do Pro incluso', 'Grupos isolados por time', 'Quotas e budgets por usuário', 'SSO + auditoria completa', 'SLA de 99.9%', 'Onboarding dedicado'],
    muted: [],
  },
]

// ── Testimonials ──────────────────────────────────────────
const nmTestimonials = [
  { quote: 'Migrei 3 produtos pro NexusMind num final de semana. Economia de R$ 2.300/mês e parou de cair quando a Anthropic engasga.', name: 'Marina V.', role: 'Tech lead · fintech SP', initial: 'MV', meta: ['ECONOMIA', 'R$ 2.300/mês'] },
  { quote: 'Setup foi o que demorou mais: 4 minutos. Troquei a base URL no Cursor e na hora estava rodando com Opus 4.7.', name: 'Caio P.', role: 'Indie hacker · POA', initial: 'CP', meta: ['SETUP', '4 MIN'] },
  { quote: 'O failover salvou meu deploy. Anthropic deu 529 quatro vezes esse mês e nenhum dos meus usuários percebeu.', name: 'Heitor M.', role: 'CTO · agência SaaS', initial: 'HM', meta: ['UPTIME', '99.94%'] },
]

// ── FAQ ───────────────────────────────────────────────────
const openFaq = ref(0)
function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? -1 : i
}
const nmFaqs = [
  { q: 'É a mesma API oficial dos providers?', a: 'Sim. O NexusMind é só uma camada de roteamento e cota — as chamadas vão direto pros endpoints oficiais da Anthropic, OpenAI e Google. Sem proxy estranho. Você usa as SDKs oficiais, só troca a base URL e a chave.' },
  { q: 'Como é cobrado? Por token ou por mês?', a: 'Por mês, com pacote de tokens incluso. Você não precisa ficar somando dólar com câmbio nem se preocupar com pico de uso surpresa. Se passar do incluso, a gente avisa antes de cobrar mais.' },
  { q: 'Em qual ferramenta consigo usar?', a: 'Em qualquer coisa que aceite uma chave OpenAI-compat: Cursor, Claude Code, Cline, Aider, Continue, OpenRouter clients, n8n, scripts em Python/Node e por aí vai. Trocou a base URL e a chave, funcionou.' },
  { q: 'Tem garantia?', a: '7 dias. Se não rolar, devolvemos 100%, sem perguntinhas chatas. Pix ou no cartão.' },
  { q: 'Posso usar a mesma chave em vários PCs?', a: 'Pode. A chave é sua. Só não compartilhe fora do seu time porque a gente monitora padrão de uso e bloqueia abuso.' },
  { q: 'Vocês são oficiais da Anthropic/OpenAI?', a: 'Não, somos um intermediário brasileiro. Operamos sob os termos de uso de cada provider e somos parceiros revenda quando a categoria existe. Empresa registrada com sede no Brasil.' },
  { q: 'Quanto demora pra liberar o acesso?', a: 'Pix confirma em segundos. Cartão em até 2 minutos. A chave já fica no painel pronta pra colar.' },
  { q: 'E se um provider cair?', a: 'Failover automático. Configure no admin a ordem de prioridade e o NexusMind redireciona pra próxima conta saudável sem você sentir.' },
]
</script>

<style scoped>
/* ── CSS custom properties ───────────────────────────────── */
.nm-home {
  --nm-bg:      #07090f;
  --nm-bg-1:    #0b0e16;
  --nm-bg-2:    #11151f;
  --nm-line:    rgba(255,255,255,0.08);
  --nm-line-2:  rgba(255,255,255,0.14);
  --nm-fg:      #e7eaf3;
  --nm-fg-1:    #b8bfd1;
  --nm-fg-2:    #8a93a8;
  --nm-fg-3:    #5b6378;
  --nm-cyan:    #5fd5ff;
  --nm-violet:  #a78bfa;
  --nm-emerald: #34d399;
  --nm-mono:    "JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, monospace;
  --nm-sans:    "Geist", ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif;
  --nm-r:       12px;

  background: var(--nm-bg);
  color: var(--nm-fg);
  font-family: var(--nm-sans);
  font-feature-settings: "ss01","cv11";
  -webkit-font-smoothing: antialiased;
  min-height: 100vh;
  overflow-x: hidden;
  position: relative;
}

/* ── Background layers ───────────────────────────────────── */
.nm-bg-grid {
  position: fixed; inset: 0; z-index: 0; pointer-events: none;
  background-image: radial-gradient(rgba(255,255,255,0.05) 1px, transparent 1px);
  background-size: 28px 28px;
  mask-image: radial-gradient(ellipse at 50% 30%, black 30%, transparent 80%);
}
.nm-bg-glow {
  position: fixed; inset: 0; z-index: 0; pointer-events: none;
  background:
    radial-gradient(700px 400px at 80% -10%, rgba(95,213,255,0.10), transparent 60%),
    radial-gradient(600px 400px at 0% 30%, rgba(167,139,250,0.08), transparent 60%);
}

/* ── Container ───────────────────────────────────────────── */
.nm-wrap { max-width: 1240px; margin: 0 auto; padding: 0 32px; position: relative; z-index: 1; }
@media (max-width: 720px) { .nm-wrap { padding: 0 20px; } }

/* ── Header ──────────────────────────────────────────────── */
.nm-header {
  position: sticky; top: 0; z-index: 50;
  backdrop-filter: blur(10px);
  background: rgba(7,9,15,0.72);
  border-bottom: 1px solid var(--nm-line);
}
.nm-header-inner { display: flex; align-items: center; justify-content: space-between; height: 64px; }
.nm-brand { display: flex; align-items: center; gap: 12px; text-decoration: none; }
.nm-brand-logo {
  width: 32px; height: 32px; border-radius: 8px; display: grid; place-items: center;
  background: #0e1320; border: 1px solid var(--nm-line-2); overflow: hidden;
}
.nm-brand-logo img { width: 28px; height: 28px; object-fit: contain; }
.nm-brand-name { display: flex; flex-direction: column; line-height: 1.05; }
.nm-brand-name b { font-weight: 600; font-size: 14px; letter-spacing: -0.01em; color: var(--nm-fg); }
.nm-brand-name span { font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); letter-spacing: 0.04em; text-transform: uppercase; }
.nm-nav { display: flex; align-items: center; gap: 28px; }
@media (max-width: 860px) { .nm-nav { display: none; } }
.nm-nav-link { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); letter-spacing: 0.04em; text-transform: uppercase; transition: color .15s; text-decoration: none; }
.nm-nav-link:hover { color: var(--nm-cyan); }
.nm-actions { display: flex; align-items: center; gap: 10px; }

/* ── Buttons ─────────────────────────────────────────────── */
.nm-btn {
  display: inline-flex; align-items: center; gap: 8px;
  height: 40px; padding: 0 16px; border-radius: 8px;
  font-family: var(--nm-mono); font-size: 12px; font-weight: 500;
  letter-spacing: 0.06em; text-transform: uppercase;
  border: 1px solid transparent; cursor: pointer; transition: all .15s;
  white-space: nowrap; text-decoration: none;
}
.nm-btn-lg { height: 48px; padding: 0 22px; font-size: 13px; }
.nm-btn-primary { background: var(--nm-cyan); color: #021018 !important; border-color: var(--nm-cyan); }
.nm-btn-primary * { color: #021018 !important; }
.nm-btn-primary:hover { background: #82e1ff; color: #021018 !important; border-color: #82e1ff; box-shadow: 0 0 0 4px rgba(95,213,255,0.15); }
.nm-btn-ghost { background: transparent; color: var(--nm-fg); border-color: var(--nm-line-2); }
.nm-btn-ghost:hover { border-color: var(--nm-cyan); color: var(--nm-cyan); }
.nm-btn-dark { background: #0e1320; color: var(--nm-fg); border-color: var(--nm-line-2); }
.nm-btn-dark:hover { background: #131a2a; }
.nm-arr { font-family: var(--nm-sans); font-weight: 500; }

/* ── Section primitives ──────────────────────────────────── */
.nm-section { padding: 96px 0; }
@media (max-width: 720px) { .nm-section { padding: 64px 0; } }
.nm-s-label {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; letter-spacing: 0.14em;
  color: var(--nm-cyan); text-transform: uppercase; margin-bottom: 16px;
}
.nm-sq { width: 6px; height: 6px; background: var(--nm-cyan); display: inline-block; }
.nm-tic { color: var(--nm-fg-3); }
.nm-s-title {
  font-size: clamp(34px, 4.2vw, 56px); font-weight: 600;
  line-height: 1.04; letter-spacing: -0.025em; margin: 0 0 18px;
  text-wrap: balance;
}
.nm-s-sub { font-size: 17px; color: var(--nm-fg-1); max-width: 640px; line-height: 1.55; margin: 0; }
.nm-accent { color: var(--nm-cyan); }
.nm-strike { text-decoration: line-through; color: var(--nm-fg-3); }

/* ── Hero ────────────────────────────────────────────────── */
.nm-hero { padding-top: 64px; padding-bottom: 80px; }
.nm-hero-grid { display: grid; grid-template-columns: 1.05fr 1fr; gap: 60px; align-items: center; }
@media (max-width: 1040px) { .nm-hero-grid { grid-template-columns: 1fr; gap: 48px; } }
.nm-eyebrow {
  display: inline-flex; align-items: center; gap: 10px;
  font-family: var(--nm-mono); font-size: 11px; text-transform: uppercase; letter-spacing: 0.14em;
  color: var(--nm-fg-1); background: #0e1320; border: 1px solid var(--nm-line-2);
  padding: 6px 12px; border-radius: 999px; margin-bottom: 18px;
}
.nm-eyebrow-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--nm-cyan); box-shadow: 0 0 12px var(--nm-cyan); }
.nm-eyebrow-v { color: var(--nm-violet); }
.nm-h1 { font-size: clamp(40px,5.8vw,76px); font-weight: 600; line-height: 0.98; letter-spacing: -0.03em; margin: 0 0 22px; text-wrap: balance; }
.nm-lede { font-size: 18px; line-height: 1.55; color: var(--nm-fg-1); max-width: 540px; margin: 0 0 32px; }
.nm-cta-row { display: flex; flex-wrap: wrap; gap: 12px; margin-bottom: 32px; }
.nm-pill-row { display: flex; flex-wrap: wrap; gap: 18px; font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); letter-spacing: 0.04em; }
.nm-pill-item { display: inline-flex; align-items: center; gap: 8px; }
.nm-dot-green { width: 8px; height: 8px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 12px rgba(52,211,153,.6); display: inline-block; }
.nm-ok { color: var(--nm-emerald); }

/* ── Console ─────────────────────────────────────────────── */
.nm-console {
  background: linear-gradient(180deg, #0d1119 0%, #090c14 100%);
  border: 1px solid var(--nm-line-2); border-radius: var(--nm-r); overflow: hidden;
  box-shadow: 0 1px 0 rgba(255,255,255,.05) inset, 0 40px 80px -30px rgba(0,0,0,.6), 0 0 0 1px rgba(95,213,255,.06);
}
.nm-console-head {
  display: flex; align-items: center; gap: 10px; height: 36px; padding: 0 14px;
  border-bottom: 1px solid var(--nm-line); background: #0a0d15;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3);
}
.nm-console-dots { display: inline-flex; gap: 6px; margin-right: 8px; }
.nm-console-dots i { width: 10px; height: 10px; border-radius: 50%; background: #1c2233; display: inline-block; }
.nm-console-title { color: var(--nm-fg-1); }
.nm-console-live { margin-left: auto; display: inline-flex; align-items: center; gap: 6px; color: var(--nm-emerald); }
.nm-live-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 10px var(--nm-emerald); display: inline-block; }
.nm-console-body { padding: 22px; }

/* SVG diagram */
.nm-diagram { width: 100%; height: 300px; display: block; }
.nm-lane { fill: none; stroke: url(#nmG1); stroke-width: 1.5; opacity: .35; stroke-dasharray: 4 6; }
.nm-lane-active { opacity: 1; stroke-dasharray: 8 4; }
.nm-lane { transition: opacity .25s ease; }
.nm-svg-label { font-family: var(--nm-mono); font-size: 10px; letter-spacing: .04em; }
.nm-svg-cyan { fill: #5fd5ff; }
.nm-svg-violet { fill: #a78bfa; }
.nm-svg-dim { fill: #5b6378; font-family: var(--nm-mono); font-size: 10px; }

.nm-console-meta { margin-top: 18px; display: grid; grid-template-columns: repeat(3,1fr); gap: 12px; font-family: var(--nm-mono); font-size: 11px; }
.nm-meta-cell { background: #0a0d15; border: 1px solid var(--nm-line); border-radius: 8px; padding: 10px 12px; }
.nm-meta-k { color: var(--nm-fg-3); letter-spacing: .08em; text-transform: uppercase; font-size: 10px; }
.nm-meta-v { color: var(--nm-fg); font-size: 16px; margin-top: 4px; letter-spacing: -.01em; font-feature-settings: "tnum"; }
.nm-meta-u { color: var(--nm-fg-3); font-size: 11px; margin-left: 2px; }

/* ── Routing diagram packets & log feed ──────────────── */
.nm-packet { fill: var(--nm-cyan); filter: drop-shadow(0 0 6px var(--nm-cyan)); opacity: 0; }
.nm-packet-1 { offset-path: path("M194 80 C 240 80, 240 160, 280 160 C 320 160, 320 70, 420 70"); animation: nm-packet-run 5.4s linear infinite; }
.nm-packet-2 { offset-path: path("M194 140 C 240 140, 240 160, 280 160 C 320 160, 320 160, 420 160"); animation: nm-packet-run 5.4s linear infinite; animation-delay: 1.8s; }
.nm-packet-3 { offset-path: path("M194 200 C 240 200, 240 160, 280 160 C 320 160, 320 250, 420 250"); animation: nm-packet-run 5.4s linear infinite; animation-delay: 3.6s; }
@keyframes nm-packet-run {
  0%   { offset-distance: 0%; opacity: 0; }
  8%   { opacity: 1; }
  42%  { opacity: 1; }
  50%, 100% { offset-distance: 100%; opacity: 0; }
}
.nm-prov-svg { fill: #0a0d15; stroke: rgba(255,255,255,0.12); transition: stroke 0.25s, fill 0.25s; }
.nm-prov-svg-active { stroke: #5fd5ff; fill: #102232; }

.nm-routing-target {
  margin-top: 14px;
  display: flex; flex-wrap: wrap; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-2);
}
.nm-rt-label { color: var(--nm-fg-3); letter-spacing: 0.06em; text-transform: uppercase; }
.nm-rt-val {
  color: var(--nm-cyan);
  background: rgba(95, 213, 255, 0.08);
  border: 1px solid rgba(95, 213, 255, 0.2);
  padding: 3px 8px;
  border-radius: 4px;
  font-feature-settings: "tnum";
}

.nm-log-feed {
  margin-top: 14px;
  background: #06080d;
  border: 1px solid var(--nm-line);
  border-radius: 8px;
  padding: 10px 14px;
  font-family: var(--nm-mono);
  font-size: 11px;
  height: 84px;
  overflow: hidden;
  position: relative;
}
.nm-log-feed::before {
  content: "$ tail -f /var/log/nexusmind/routes.log";
  display: block;
  font-size: 10px; color: var(--nm-fg-3);
  letter-spacing: 0.06em;
  padding-bottom: 6px;
  border-bottom: 1px dashed var(--nm-line);
  margin-bottom: 6px;
}
.nm-log-line {
  display: flex; gap: 8px; align-items: baseline;
  line-height: 1.5;
  white-space: nowrap; overflow: hidden;
}
.nm-log-t { color: var(--nm-fg-3); }
.nm-log-method { color: var(--nm-violet); }
.nm-log-arr { color: var(--nm-fg-3); }
.nm-log-prov { color: var(--nm-cyan); }
.nm-log-model { color: var(--nm-fg-1); }
.nm-log-latency { color: var(--nm-fg-2); margin-left: auto; }
.nm-log-ok { color: var(--nm-emerald); }

/* ── Tutorial section ───────────────────────────────────── */
.nm-tut-grid { display: grid; grid-template-columns: 1.6fr 1fr; gap: 36px; margin-top: 48px; }
@media (max-width: 920px) { .nm-tut-grid { grid-template-columns: 1fr; gap: 28px; } }
.nm-video-wrap {
  position: relative;
  aspect-ratio: 16 / 9;
  width: 100%;
  border-radius: var(--nm-r);
  overflow: hidden;
  background: #0a0d15;
  border: 1px solid var(--nm-line-2);
  box-shadow: 0 1px 0 rgba(255,255,255,0.05) inset, 0 30px 60px -30px rgba(0,0,0,0.7);
}
.nm-video-wrap iframe { position: absolute; inset: 0; width: 100%; height: 100%; border: 0; }
.nm-video-placeholder {
  position: absolute; inset: 0;
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 18px;
  background:
    radial-gradient(600px 300px at 50% 60%, rgba(95,213,255,0.08), transparent 60%),
    repeating-linear-gradient(45deg, #0a0d15 0px, #0a0d15 12px, #0c1019 12px, #0c1019 24px);
  cursor: pointer;
  transition: background 0.2s;
}
.nm-video-placeholder:hover { background:
  radial-gradient(700px 350px at 50% 60%, rgba(95,213,255,0.14), transparent 60%),
  repeating-linear-gradient(45deg, #0a0d15 0px, #0a0d15 12px, #0c1019 12px, #0c1019 24px); }
.nm-video-badge {
  position: absolute; top: 18px; left: 18px;
  display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase;
  color: var(--nm-cyan);
  background: rgba(7,9,15,0.85); border: 1px solid var(--nm-line-2);
  padding: 6px 12px; border-radius: 999px;
}
.nm-video-time {
  position: absolute; top: 18px; right: 18px;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-2);
  background: rgba(7,9,15,0.85); border: 1px solid var(--nm-line-2);
  padding: 6px 12px; border-radius: 999px;
}
.nm-play-btn {
  width: 82px; height: 82px; border-radius: 50%;
  background: var(--nm-cyan); color: #021018;
  display: grid; place-items: center; border: none; cursor: pointer;
  transition: transform 0.15s;
}
.nm-play-btn:hover { transform: scale(1.06); }
.nm-play-btn svg { margin-left: 6px; }
@keyframes nm-play-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(95, 213, 255, 0.5); }
  50%      { box-shadow: 0 0 0 20px rgba(95, 213, 255, 0); }
}
.nm-video-hint { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); letter-spacing: 0.08em; text-transform: uppercase; text-align: center; padding: 0 24px; }
.nm-video-hint code { color: var(--nm-cyan); }
.nm-video-strip {
  position: absolute; bottom: 16px; left: 18px; right: 18px;
  display: flex; gap: 14px; align-items: center; justify-content: space-between;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3);
}
.nm-vs-left { display: inline-flex; align-items: center; gap: 8px; }
.nm-vs-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--nm-cyan); display: inline-block; }
.nm-vs-progress { flex: 1; max-width: 200px; height: 3px; background: rgba(255,255,255,0.1); border-radius: 2px; overflow: hidden; position: relative; }
.nm-vs-progress::after { content: ""; display: block; width: 22%; height: 100%; background: var(--nm-cyan); }

.nm-tut-h { font-family: var(--nm-mono); font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: var(--nm-fg-3); margin: 0 0 16px; }
.nm-tut-list { list-style: none; margin: 0; padding: 0; }
.nm-tut-list li {
  display: grid; grid-template-columns: 56px 1fr;
  align-items: baseline; gap: 14px;
  padding: 14px 0; border-top: 1px solid var(--nm-line);
  font-size: 15px; color: var(--nm-fg);
  transition: padding-left 0.15s, color 0.15s;
  cursor: pointer;
}
.nm-tut-list li:hover { padding-left: 8px; color: var(--nm-cyan); }
.nm-tut-list li:last-child { border-bottom: 1px solid var(--nm-line); }
.nm-tut-time { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-cyan); letter-spacing: 0.04em; font-feature-settings: "tnum"; }

/* ── Integrations grid ─────────────────────────────────── */
.nm-int-grid { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; margin-top: 48px; }
@media (max-width: 1040px) { .nm-int-grid { grid-template-columns: repeat(2,1fr); } }
@media (max-width: 540px)  { .nm-int-grid { grid-template-columns: 1fr; } }
.nm-int-card {
  background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r);
  padding: 20px; position: relative;
  transition: border-color 0.2s, transform 0.2s, background 0.2s;
}
.nm-int-card:hover { border-color: var(--nm-cyan); background: linear-gradient(180deg, #0e1923 0%, #0a0d15 100%); transform: translateY(-2px); }
.nm-int-row { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.nm-int-icon { width: 36px; height: 36px; background: #06080d; border: 1px solid var(--nm-line-2); border-radius: 8px; display: grid; place-items: center; font-family: var(--nm-mono); font-size: 16px; font-weight: 600; color: var(--nm-cyan); }
.nm-int-name { font-size: 16px; font-weight: 600; letter-spacing: -0.01em; }
.nm-int-where { font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); letter-spacing: 0.04em; text-transform: uppercase; margin-bottom: 10px; }
.nm-int-env { background: #06080d; border: 1px solid var(--nm-line); font-family: var(--nm-mono); font-size: 11px; padding: 8px 10px; border-radius: 6px; color: var(--nm-fg-2); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.nm-int-env b { color: var(--nm-cyan); font-weight: 500; }
.nm-int-url { color: var(--nm-violet); }
.nm-int-tested {
  position: absolute; top: 18px; right: 18px;
  font-family: var(--nm-mono); font-size: 9px; color: var(--nm-emerald);
  display: inline-flex; align-items: center; gap: 4px;
  letter-spacing: 0.1em; text-transform: uppercase;
}
.nm-int-d { width: 5px; height: 5px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 6px var(--nm-emerald); display: inline-block; }
.nm-int-more {
  margin-top: 32px; padding: 20px 24px;
  border: 1px dashed var(--nm-line-2); border-radius: var(--nm-r);
  font-family: var(--nm-mono); font-size: 12px;
  color: var(--nm-fg-2); letter-spacing: 0.04em; line-height: 1.6;
}
.nm-int-more b { color: var(--nm-cyan); font-weight: 500; }

/* ── Floating CTA pill ───────────────────────────────────── */
.nm-float-cta {
  position: fixed; left: 50%; bottom: 20px;
  transform: translateX(-50%) translateY(120%);
  z-index: 80;
  transition: transform 0.35s cubic-bezier(0.5, 1.5, 0.3, 1);
  pointer-events: none;
}
.nm-float-visible { transform: translateX(-50%) translateY(0); pointer-events: auto; }
.nm-float-inner {
  display: flex; align-items: center; gap: 16px;
  padding: 8px 8px 8px 18px;
  background: rgba(11, 14, 22, 0.92);
  backdrop-filter: blur(16px);
  border: 1px solid var(--nm-line-2);
  border-radius: 999px;
  box-shadow: 0 24px 60px -20px rgba(0,0,0,0.7);
  font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-1);
}
.nm-float-prompt { color: var(--nm-cyan); font-weight: 600; letter-spacing: 0.02em; }
.nm-float-sep { color: var(--nm-line-2); margin: 0 -8px; }
.nm-float-link {
  color: var(--nm-fg-2); letter-spacing: 0.04em; text-transform: uppercase; font-size: 11px;
  transition: color 0.15s; text-decoration: none;
}
.nm-float-link:hover { color: var(--nm-cyan); }
.nm-float-btn {
  height: 36px; padding: 0 16px;
  background: var(--nm-cyan); color: #021018 !important;
  border: none; border-radius: 999px;
  font-family: var(--nm-mono); font-size: 12px; font-weight: 500;
  letter-spacing: 0.04em; text-transform: uppercase;
  display: inline-flex; align-items: center; gap: 6px;
  cursor: pointer; text-decoration: none;
  transition: background 0.15s;
}
.nm-float-btn * { color: #021018 !important; }
.nm-float-btn:hover { background: #82e1ff; color: #021018 !important; }
.nm-float-close {
  width: 36px; height: 36px;
  display: grid; place-items: center;
  background: transparent; border: none;
  color: var(--nm-fg-3); cursor: pointer;
  border-radius: 50%;
  font-family: var(--nm-sans); font-size: 16px;
}
.nm-float-close:hover { color: var(--nm-fg); background: rgba(255,255,255,0.05); }
@media (max-width: 720px) {
  .nm-float-link, .nm-float-sep { display: none; }
}

/* ── Scroll reveal ───────────────────────────────────────── */
.nm-reveal {
  opacity: 0;
  transform: translateY(20px);
  transition: opacity 0.6s ease, transform 0.6s ease;
}
.nm-reveal-in { opacity: 1; transform: translateY(0); }
@media (prefers-reduced-motion: reduce) {
  .nm-reveal, .nm-reveal-in { opacity: 1; transform: none; transition: none; }
  .nm-packet { animation: none; }
}

/* ── Steps ───────────────────────────────────────────────── */
.nm-steps { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; }
@media (max-width: 920px) { .nm-steps { grid-template-columns: 1fr; } }
.nm-step { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 28px; transition: border-color .2s; }
.nm-step:hover { border-color: var(--nm-line-2); }
.nm-step-num { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-3); letter-spacing: .1em; }
.nm-step-num b { color: var(--nm-cyan); font-weight: 500; }
.nm-step h3 { font-size: 22px; font-weight: 600; letter-spacing: -.015em; margin: 14px 0 8px; }
.nm-step p { color: var(--nm-fg-1); font-size: 15px; line-height: 1.55; margin: 0 0 18px; }
.nm-codeblock { background: #06080d; border: 1px solid var(--nm-line); font-family: var(--nm-mono); font-size: 12px; padding: 12px 14px; border-radius: 8px; color: var(--nm-fg-1); overflow: hidden; }
.nm-c { color: var(--nm-fg-3); }
.nm-k { color: var(--nm-cyan); }
.nm-v { color: var(--nm-violet); }
.nm-s { color: var(--nm-emerald); }

/* ── Math ────────────────────────────────────────────────── */
.nm-math-grid { display: grid; grid-template-columns: 1.1fr 1fr; gap: 60px; align-items: center; margin-top: 56px; }
@media (max-width: 1000px) { .nm-math-grid { grid-template-columns: 1fr; gap: 40px; } }
.nm-compare { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); overflow: hidden; }
.nm-cmp-row { display: grid; grid-template-columns: 1.2fr 1fr 1fr; border-bottom: 1px solid var(--nm-line); font-family: var(--nm-mono); font-size: 13px; }
.nm-cmp-row:last-child { border-bottom: none; }
.nm-cmp-header { background: #0a0d15; }
.nm-cmp-header .nm-cmp-cell { font-size: 10px; letter-spacing: .1em; text-transform: uppercase; color: var(--nm-fg-3); padding: 14px 18px; }
.nm-cmp-header .nm-cmp-ours { color: var(--nm-cyan); }
.nm-cmp-cell { padding: 16px 18px; border-right: 1px solid var(--nm-line); color: var(--nm-fg-1); }
.nm-cmp-cell:last-child { border-right: none; }
.nm-cmp-label { color: var(--nm-fg-2); }
.nm-cmp-ours { color: var(--nm-fg); font-weight: 500; background: rgba(95,213,255,.03); }
.nm-cmp-total { background: #0a0d15; }
.nm-cmp-total .nm-cmp-cell { font-size: 15px; color: var(--nm-fg); padding: 18px; font-weight: 600; }
.nm-cmp-total .nm-cmp-ours { color: var(--nm-cyan); }
.nm-math-points { display: flex; flex-direction: column; gap: 24px; }
.nm-math-pt { border-left: 2px solid var(--nm-line-2); padding: 4px 0 4px 20px; cursor: default; }
.nm-math-pt-active { border-left-color: var(--nm-cyan); }
.nm-pt-tag { font-family: var(--nm-mono); font-size: 11px; letter-spacing: .1em; color: var(--nm-violet); text-transform: uppercase; }
.nm-math-pt h4 { font-size: 18px; font-weight: 600; letter-spacing: -.01em; margin: 6px 0 8px; }
.nm-math-pt p { color: var(--nm-fg-1); font-size: 14.5px; line-height: 1.55; margin: 0; }

/* ── Providers ───────────────────────────────────────────── */
.nm-providers { display: grid; grid-template-columns: repeat(3,1fr); gap: 16px; margin-top: 48px; }
@media (max-width: 1080px) { .nm-providers { grid-template-columns: repeat(2,1fr); } }
@media (max-width: 720px) { .nm-providers { grid-template-columns: 1fr; } }
.nm-prov-card { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 22px; position: relative; transition: all .2s; }
.nm-prov-card:hover { border-color: var(--nm-line-2); transform: translateY(-2px); }
.nm-prov-status { position: absolute; top: 18px; right: 18px; font-family: var(--nm-mono); font-size: 10px; color: var(--nm-emerald); display: inline-flex; align-items: center; gap: 6px; }
.nm-prov-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 8px var(--nm-emerald); display: inline-block; }
.nm-prov-ico { width: 40px; height: 40px; border-radius: 10px; display: grid; place-items: center; font-family: var(--nm-mono); font-size: 16px; font-weight: 600; background: #0a0d15; border: 1px solid; margin-bottom: 14px; }
.nm-prov-card h4 { font-size: 16px; font-weight: 600; margin: 0 0 6px; letter-spacing: -.01em; }
.nm-prov-models { display: flex; flex-wrap: wrap; gap: 6px; font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-2); letter-spacing: .04em; line-height: 1.6; }
.nm-prov-models b { flex-basis: 100%; color: var(--nm-fg); font-weight: 500; margin-bottom: 2px; }
.nm-prov-models span { display: inline-block; padding: 3px 7px; border-radius: 999px; background: rgba(255,255,255,.04); border: 1px solid rgba(255,255,255,.06); color: var(--nm-fg-1); }
.nm-model-note { margin-top: 28px; padding: 18px 22px; border: 1px dashed var(--nm-line-2); border-radius: var(--nm-r); color: var(--nm-fg-2); font-family: var(--nm-mono); font-size: 12px; line-height: 1.6; }
.nm-model-note b { color: var(--nm-cyan); font-weight: 500; }

/* ── Pricing ─────────────────────────────────────────────── */
.nm-pricing-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; align-items: stretch; }
@media (max-width: 960px) { .nm-pricing-grid { grid-template-columns: 1fr; } }
.nm-plan { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 32px; position: relative; display: flex; flex-direction: column; }
.nm-plan-feat { border-color: var(--nm-cyan); background: linear-gradient(180deg, #0e1923 0%, #0a0d15 100%); }
.nm-plan-badge { position: absolute; top: -10px; left: 32px; font-family: var(--nm-mono); font-size: 10px; letter-spacing: .14em; text-transform: uppercase; background: var(--nm-cyan); color: #021018; padding: 4px 10px; border-radius: 4px; }
.nm-plan-type { font-family: var(--nm-mono); font-size: 11px; letter-spacing: .14em; text-transform: uppercase; color: var(--nm-fg-2); margin-bottom: 12px; }
.nm-plan h3 { font-size: 24px; font-weight: 600; letter-spacing: -.015em; margin: 0 0 8px; }
.nm-plan-desc { font-size: 14px; color: var(--nm-fg-2); line-height: 1.5; margin: 0 0 22px; min-height: 42px; }
.nm-plan-price { display: flex; align-items: baseline; gap: 4px; margin-bottom: 6px; }
.nm-price-cur { font-family: var(--nm-mono); font-size: 14px; color: var(--nm-fg-2); }
.nm-price-amt { font-size: 56px; font-weight: 600; letter-spacing: -.04em; line-height: 1; }
.nm-price-per { font-family: var(--nm-mono); font-size: 13px; color: var(--nm-fg-3); margin-left: 4px; }
.nm-plan-old { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-3); text-decoration: line-through; margin-bottom: 22px; }
.nm-plan-old-spacer { text-decoration: none; visibility: hidden; }
.nm-plan-features { list-style: none; margin: 0 0 24px; padding: 0; flex: 1; }
.nm-plan-features li { font-size: 14px; color: var(--nm-fg-1); padding: 8px 0 8px 24px; position: relative; border-bottom: 1px solid var(--nm-line); }
.nm-plan-features li:last-child { border-bottom: none; }
.nm-plan-features li::before { content: "→"; position: absolute; left: 0; top: 8px; font-family: var(--nm-mono); color: var(--nm-cyan); font-size: 14px; }
.nm-plan-features li.nm-feat-muted { color: var(--nm-fg-3); }
.nm-plan-features li.nm-feat-muted::before { content: "·"; color: var(--nm-fg-3); }

/* ── Proof ───────────────────────────────────────────────── */
.nm-proof-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; }
@media (max-width: 1000px) { .nm-proof-grid { grid-template-columns: 1fr; } }
.nm-tcard { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 28px; display: flex; flex-direction: column; }
.nm-tcard-meta { font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); margin-bottom: 16px; letter-spacing: .08em; text-transform: uppercase; display: flex; gap: 8px; }
.nm-t-sep { color: var(--nm-line-2); }
.nm-t-ec { color: var(--nm-emerald); }
.nm-tcard-quote { font-size: 17px; line-height: 1.5; letter-spacing: -.01em; color: var(--nm-fg); margin: 0 0 24px; text-wrap: pretty; flex: 1; }
.nm-tcard-who { display: flex; align-items: center; gap: 12px; }
.nm-tcard-av { width: 36px; height: 36px; border-radius: 50%; background: linear-gradient(135deg, var(--nm-cyan), var(--nm-violet)); display: grid; place-items: center; font-family: var(--nm-mono); font-size: 13px; font-weight: 600; color: #021018; flex-shrink: 0; }
.nm-tcard-who b { display: block; font-size: 14px; font-weight: 500; }
.nm-tcard-who span { display: block; font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); }

/* ── FAQ ─────────────────────────────────────────────────── */
.nm-faq-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0 36px; margin-top: 32px; }
@media (max-width: 880px) { .nm-faq-grid { grid-template-columns: 1fr; } }
.nm-faq-item { border-top: 1px solid var(--nm-line); }
.nm-faq-item button {
  display: flex; align-items: center; justify-content: space-between; gap: 18px;
  width: 100%; padding: 20px 0; background: transparent; border: none; cursor: pointer; text-align: left;
  color: var(--nm-fg); font-size: 15.5px; font-weight: 500; letter-spacing: -.01em; font-family: inherit;
}
.nm-faq-item button:hover { color: var(--nm-cyan); }
.nm-faq-num { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-3); margin-right: 14px; }
.nm-faq-plus { color: var(--nm-cyan); font-family: var(--nm-mono); font-size: 18px; transition: transform .2s; flex-shrink: 0; }
.nm-faq-open .nm-faq-plus { transform: rotate(45deg); }
.nm-faq-answer { max-height: 0; overflow: hidden; transition: max-height .25s ease, padding .25s ease; }
.nm-faq-open .nm-faq-answer { max-height: 400px; padding: 0 0 24px 30px; }
.nm-faq-answer p { color: var(--nm-fg-1); font-size: 14.5px; line-height: 1.6; margin: 0; max-width: 56ch; }

/* ── Final CTA ───────────────────────────────────────────── */
.nm-final-cta {
  background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r);
  padding: 56px; text-align: center; position: relative; overflow: hidden;
}
.nm-final-cta::before {
  content: ""; position: absolute; inset: 0; pointer-events: none;
  background: radial-gradient(600px 300px at 50% 0%, rgba(95,213,255,.12), transparent 70%);
}
.nm-final-cta .nm-s-label { display: inline-flex; }
.nm-final-cta .nm-s-title { font-size: clamp(34px,4vw,52px); position: relative; }
.nm-cta-p { font-size: 17px; color: var(--nm-fg-1); margin: 0 auto 28px; max-width: 520px; position: relative; }
.nm-cta-btns { display: inline-flex; gap: 12px; position: relative; flex-wrap: wrap; justify-content: center; }
@media (max-width: 720px) { .nm-final-cta { padding: 40px 24px; } }

/* ── Footer ──────────────────────────────────────────────── */
.nm-footer { border-top: 1px solid var(--nm-line); margin-top: 80px; padding: 48px 0 28px; position: relative; z-index: 1; }
.nm-footer-grid { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr; gap: 36px; }
@media (max-width: 880px) { .nm-footer-grid { grid-template-columns: 1fr 1fr; } }
.nm-footer-col h5 { font-family: var(--nm-mono); font-size: 11px; letter-spacing: .14em; text-transform: uppercase; color: var(--nm-fg-3); margin: 0 0 16px; }
.nm-footer-col a { display: block; font-size: 14px; color: var(--nm-fg-1); padding: 6px 0; text-decoration: none; }
.nm-footer-col a:hover { color: var(--nm-cyan); }
.nm-footer-desc { font-size: 14px; color: var(--nm-fg-2); line-height: 1.55; max-width: 360px; margin: 12px 0 0; }
.nm-footer-bot {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: 36px; padding-top: 24px; border-top: 1px solid var(--nm-line);
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); letter-spacing: .04em;
}
@media (max-width: 720px) { .nm-footer-bot { flex-direction: column; gap: 12px; text-align: center; } }
</style>
