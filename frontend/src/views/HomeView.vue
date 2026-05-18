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

  <!-- Route-Cortexx Redesign -->
  <div v-else class="rc-home">
    <div class="rc-bg-grid"></div>
    <div class="rc-bg-glow"></div>

    <!-- ── Header ─────────────────────────────────────────── -->
    <header class="rc-header">
      <div class="rc-wrap rc-header-inner">
        <a href="/home" class="rc-brand">
          <div class="rc-brand-logo">
            <img :src="siteLogo || '/logo.png'" alt="Route-Cortexx" />
          </div>
          <div class="rc-brand-name">
            <b>Route-Cortexx</b>
            <span>AI ROUTING · BR</span>
          </div>
        </a>

        <nav class="rc-nav">
          <a href="#rc-how" class="rc-nav-link">$ como_funciona</a>
          <a href="#rc-math" class="rc-nav-link">$ matemática</a>
          <a href="#rc-pricing" class="rc-nav-link">$ preço</a>
          <a href="#rc-faq" class="rc-nav-link">$ faq</a>
        </nav>

        <div class="rc-actions">
          <LocaleSwitcher />
          <router-link v-if="isAuthenticated" :to="dashboardPath" class="rc-btn rc-btn-ghost">
            {{ t('home.dashboard') }}
          </router-link>
          <router-link v-else to="/login" class="rc-btn rc-btn-ghost">
            {{ t('home.login') }}
          </router-link>
          <a href="#rc-pricing" class="rc-btn rc-btn-primary">
            Começar <span class="rc-arr">→</span>
          </a>
        </div>
      </div>
    </header>

    <!-- ── Main ───────────────────────────────────────────── -->
    <main>

      <!-- Hero -->
      <section id="rc-top" class="rc-hero rc-wrap">
        <div class="rc-hero-grid">
          <div>
            <div class="rc-eyebrow">
              <span class="rc-eyebrow-dot"></span>
              <span>Route-Cortexx · <span class="rc-eyebrow-v">v2.1</span> · BR</span>
            </div>
            <h1 class="rc-h1">
              Uma chave.<br />
              <span class="rc-accent">Todos os modelos</span><br />
              de IA.
            </h1>
            <p class="rc-lede">
              Roteamento inteligente para Claude, GPT e Gemini —
              no real, com Pix. Painel pra monitorar tokens, custo
              e failover automático quando um provider engasga.
            </p>
            <div class="rc-cta-row">
              <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="rc-btn rc-btn-primary rc-btn-lg">
                {{ isAuthenticated ? 'Abrir painel' : 'Quero minha chave' }}
                <span class="rc-arr">→</span>
              </router-link>
              <a href="#rc-how" class="rc-btn rc-btn-ghost rc-btn-lg">Como funciona</a>
            </div>
            <div class="rc-pill-row">
              <span class="rc-pill-item">
                <span class="rc-dot-green"></span> 3 providers · 14 modelos
              </span>
              <span class="rc-pill-item">
                <span class="rc-ok">●</span> 99.94% uptime · 30d
              </span>
              <span class="rc-pill-item">★ 4.9 · 100+ reviews</span>
            </div>
          </div>

          <!-- Routing console -->
          <div class="rc-console">
            <div class="rc-console-head">
              <span class="rc-console-dots"><i /><i /><i /></span>
              <span class="rc-console-title">router.cortexx.online</span>
              <span class="rc-console-live">
                <i class="rc-live-dot"></i> ROUTING
              </span>
            </div>
            <div class="rc-console-body">
              <!-- Routing SVG diagram (v2: with packets + active highlight) -->
              <svg class="rc-diagram" viewBox="0 0 560 320" preserveAspectRatio="xMidYMid meet">
                <defs>
                  <linearGradient id="rcG1" x1="0" y1="0" x2="1" y2="0">
                    <stop offset="0%" stop-color="#5fd5ff" />
                    <stop offset="100%" stop-color="#a78bfa" />
                  </linearGradient>
                  <radialGradient id="rcCoreGlow" cx="50%" cy="50%" r="50%">
                    <stop offset="0%" stop-color="rgba(95,213,255,0.45)" />
                    <stop offset="100%" stop-color="rgba(95,213,255,0)" />
                  </radialGradient>
                </defs>
                <g opacity="0.07">
                  <line v-for="n in 7" :key="n" x1="0" :y1="n*46+20" x2="560" :y2="n*46+20" stroke="#5fd5ff" stroke-dasharray="2 6" />
                </g>
                <text x="14" y="22" class="rc-svg-dim">$ ingress</text>
                <text x="540" y="22" class="rc-svg-dim" text-anchor="end">$ providers</text>

                <!-- Request nodes -->
                <g v-for="(req, i) in rcRequests" :key="i" :opacity="activeRoute === i ? 1 : 0.55">
                  <rect :x="14" :y="req.y-12" width="180" height="24" rx="4" fill="#0a0d15" :stroke="activeRoute === i ? '#5fd5ff' : 'rgba(255,255,255,0.11)'" />
                  <circle cx="26" :cy="req.y" r="3" :fill="activeRoute === i ? '#5fd5ff' : '#a78bfa'" />
                  <text x="36" :y="req.y+4" :class="['rc-svg-label', activeRoute === i ? 'rc-svg-cyan' : 'rc-svg-dim']">{{ req.label }}</text>
                </g>

                <!-- Lanes to core -->
                <path
                  v-for="(req, i) in rcRequests"
                  :key="'l'+i"
                  :d="`M194 ${req.y} C 240 ${req.y}, 240 160, 280 160`"
                  stroke="url(#rcG1)"
                  :class="['rc-lane', activeRoute === i ? 'rc-lane-active' : '']"
                />

                <!-- Core (router) -->
                <circle cx="280" cy="160" r="58" fill="url(#rcCoreGlow)" />
                <rect x="240" y="138" width="80" height="44" rx="6" fill="#0d1119" stroke="#5fd5ff" />
                <text x="280" y="155" class="rc-svg-cyan" text-anchor="middle" font-size="11" font-weight="600" font-family="JetBrains Mono,monospace" style="letter-spacing: 0.06em">ROUTER</text>
                <text x="280" y="170" class="rc-svg-dim" text-anchor="middle">cortexx.online</text>

                <!-- Provider lanes + boxes -->
                <g v-for="(prov, i) in rcProviders" :key="'p'+i">
                  <path
                    :d="`M320 160 C 380 160, 380 ${prov.y}, 420 ${prov.y}`"
                    stroke="url(#rcG1)"
                    :class="['rc-lane', activeRoute === i ? 'rc-lane-active' : '']"
                  />
                  <rect :x="420" :y="prov.y-12" width="126" height="24" rx="4" :class="['rc-prov-svg', activeRoute === i ? 'rc-prov-svg-active' : '']" />
                  <rect :x="426" :y="prov.y-7" width="14" height="14" rx="2" :fill="activeRoute === i ? '#5fd5ff' : '#1c2233'" />
                  <text :x="448" :y="prov.y+4" :class="['rc-svg-label', activeRoute === i ? 'rc-svg-cyan' : 'rc-svg-dim']">{{ prov.name }}</text>
                  <text x="540" :y="prov.y+4" class="rc-svg-dim" text-anchor="end" font-size="9">{{ activeRoute === i ? '200 OK' : '—' }}</text>
                </g>

                <!-- Traveling packet (one per active lane, key forces re-mount → restart animation) -->
                <circle v-if="activeRoute === 0" :key="'pk0-'+packetTick" class="rc-packet rc-packet-1" r="4" />
                <circle v-if="activeRoute === 1" :key="'pk1-'+packetTick" class="rc-packet rc-packet-2" r="4" />
                <circle v-if="activeRoute === 2" :key="'pk2-'+packetTick" class="rc-packet rc-packet-3" r="4" />
              </svg>

              <!-- Routing target tag -->
              <div class="rc-routing-target">
                <span class="rc-rt-label">$ now routing</span>
                <span class="rc-rt-val">{{ rcRequests[activeRoute].method }}</span>
                <span class="rc-rt-label">→</span>
                <span class="rc-rt-val">{{ rcRequests[activeRoute].prov }}</span>
                <span class="rc-rt-val">{{ rcRequests[activeRoute].model }}</span>
              </div>

              <!-- Live log feed -->
              <div class="rc-log-feed">
                <div v-for="(log, i) in logs" :key="log.id" class="rc-log-line" :style="{ opacity: 1 - i * 0.32 }">
                  <span class="rc-log-t">[{{ log.t }}]</span>
                  <span class="rc-log-method">{{ log.path }}</span>
                  <span class="rc-log-arr">→</span>
                  <span class="rc-log-prov">{{ log.prov }}</span>
                  <span class="rc-log-model">{{ log.model }}</span>
                  <span class="rc-log-latency">{{ log.latency }}ms</span>
                  <span class="rc-log-ok">200</span>
                </div>
              </div>

              <div class="rc-console-meta">
                <div class="rc-meta-cell">
                  <div class="rc-meta-k">req/s</div>
                  <div class="rc-meta-v">428<span class="rc-meta-u">/s</span></div>
                </div>
                <div class="rc-meta-cell">
                  <div class="rc-meta-k">latência p95</div>
                  <div class="rc-meta-v">312<span class="rc-meta-u">ms</span></div>
                </div>
                <div class="rc-meta-cell">
                  <div class="rc-meta-k">economia média</div>
                  <div class="rc-meta-v">R$ 851<span class="rc-meta-u">/mês</span></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ── Tutorial ─────────────────────────────────────── -->
      <section id="rc-tutorial" class="rc-section rc-wrap">
        <div class="rc-s-label rc-reveal"><span class="rc-sq"></span> $ tutorial <span class="rc-tic">--duration=90s</span></div>
        <h2 class="rc-s-title rc-reveal">Veja em 90 segundos.</h2>
        <p class="rc-s-sub rc-reveal">Do Pix até a primeira chamada. Sem áudio, sem firula — só os cliques na tela.</p>

        <div class="rc-tut-grid">
          <div class="rc-reveal">
            <div class="rc-video-wrap">
              <!-- ⬇️ COLE O LINK DO VÍDEO AQUI: substitua videoEmbedUrl (no <script>) por uma URL embed do YouTube
                   ex.: "https://www.youtube.com/embed/SEU_VIDEO_ID" -->
              <iframe
                v-if="videoEmbedUrl"
                :src="videoEmbedUrl"
                title="Tutorial Route-Cortexx"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowfullscreen
              ></iframe>
              <div v-else class="rc-video-placeholder">
                <div class="rc-video-badge">✦ TUTORIAL · 90s</div>
                <div class="rc-video-time">SEM ÁUDIO</div>
                <button class="rc-play-btn" aria-label="Reproduzir tutorial">
                  <svg viewBox="0 0 24 24" width="32" height="32" fill="currentColor">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </button>
                <div class="rc-video-hint">Cole o link do vídeo em <code>videoEmbedUrl</code></div>
                <div class="rc-video-strip">
                  <span class="rc-vs-left"><span class="rc-vs-dot"></span> ROUTE-CORTEXX · /tutorial</span>
                  <span class="rc-vs-progress"></span>
                  <span>00:22 / 01:30</span>
                </div>
              </div>
            </div>
          </div>

          <div class="rc-tut-chapters rc-reveal">
            <h4 class="rc-tut-h">$ capítulos --chapters</h4>
            <ul class="rc-tut-list">
              <li v-for="(c, i) in rcChapters" :key="i">
                <span class="rc-tut-time">{{ c.t }}</span>
                <span>{{ c.txt }}</span>
              </li>
            </ul>
          </div>
        </div>
      </section>

      <!-- Como funciona -->
      <section id="rc-how" class="rc-section rc-wrap">
        <div class="rc-s-label"><span class="rc-sq"></span> $ como_funciona <span class="rc-tic">--steps=3</span></div>
        <h2 class="rc-s-title">Em 3 passos.<br />Sem firula.</h2>
        <p class="rc-s-sub">Você não precisa migrar nada, criar conta em provider gringo, nem pagar dólar. Só trocar a base URL.</p>

        <div class="rc-steps">
          <article class="rc-step">
            <div class="rc-step-num"><b>01</b> /// CRIAR CHAVE</div>
            <h3>Pagar e receber a chave</h3>
            <p>Pix confirma na hora. A chave aparece no seu painel, no mesmo email da compra.</p>
            <div class="rc-codeblock">
              <span class="rc-c"># painel.cortexx.online</span><br />
              <span class="rc-k">sk_live</span>_<span class="rc-v">7c9f...8b2a</span>
            </div>
          </article>
          <article class="rc-step">
            <div class="rc-step-num"><b>02</b> /// CONFIGURAR</div>
            <h3>Trocar a base URL</h3>
            <p>Cole a chave e a base URL no seu Cursor, Claude Code, script ou SDK. Funciona em qualquer client OpenAI-compat.</p>
            <div class="rc-codeblock">
              <span class="rc-k">ANTHROPIC_BASE_URL</span>=<br />
              <span class="rc-s">"https://router.cortexx.online"</span>
            </div>
          </article>
          <article class="rc-step">
            <div class="rc-step-num"><b>03</b> /// USAR</div>
            <h3>Pedir e roteamos</h3>
            <p>Seu request chega no router, escolhemos o provider saudável mais barato pro modelo pedido e devolvemos a resposta.</p>
            <div class="rc-codeblock">
              <span class="rc-c">→ model=claude-opus-4-7</span><br />
              <span class="rc-c">→ routed to provider AN_03</span><br />
              <span class="rc-s">← 200 OK · 482ms</span>
            </div>
          </article>
        </div>
      </section>

      <!-- Matemática -->
      <section id="rc-math" class="rc-section rc-wrap">
        <div class="rc-s-label"><span class="rc-sq"></span> $ matemática <span class="rc-tic">--why=cheaper</span></div>
        <h2 class="rc-s-title">Por que custa <span class="rc-accent">R$ 249</span><br />no lugar de R$ 1.020?</h2>
        <p class="rc-s-sub">Sem pegadinha. Mesmos modelos, mesma API, mesma performance. A diferença vem de três fatores.</p>

        <div class="rc-math-grid">
          <div class="rc-compare">
            <div class="rc-cmp-row rc-cmp-header">
              <div class="rc-cmp-cell">Assinatura mensal</div>
              <div class="rc-cmp-cell">Direto no provider</div>
              <div class="rc-cmp-cell rc-cmp-ours">Via Route-Cortexx</div>
            </div>
            <div v-for="row in rcMathRows" :key="row.label" class="rc-cmp-row">
              <div class="rc-cmp-cell rc-cmp-label">{{ row.label }}</div>
              <div class="rc-cmp-cell"><span class="rc-strike">{{ row.direct }}</span></div>
              <div class="rc-cmp-cell rc-cmp-ours">incluso</div>
            </div>
            <div class="rc-cmp-row rc-cmp-total">
              <div class="rc-cmp-cell">Total / mês</div>
              <div class="rc-cmp-cell"><span class="rc-strike">R$ 1.020</span></div>
              <div class="rc-cmp-cell rc-cmp-ours">R$ 249</div>
            </div>
          </div>

          <div class="rc-math-points">
            <div
              v-for="(pt, i) in rcMathPoints"
              :key="i"
              class="rc-math-pt"
              :class="{ 'rc-math-pt-active': activeMath === i }"
              @mouseenter="activeMath = i"
            >
              <div class="rc-pt-tag">{{ pt.tag }}</div>
              <h4>{{ pt.title }}</h4>
              <p>{{ pt.body }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ── Integrações ──────────────────────────────────── -->
      <section id="rc-integrations" class="rc-section rc-wrap">
        <div class="rc-s-label rc-reveal"><span class="rc-sq"></span> $ integrações <span class="rc-tic">--openai-compat=true</span></div>
        <h2 class="rc-s-title rc-reveal">Cola a key onde<br />você já trabalha.</h2>
        <p class="rc-s-sub rc-reveal">Trocou a base URL e a chave, funcionou. Qualquer client OpenAI-compat funciona — sem migração, sem refactor.</p>

        <div class="rc-int-grid">
          <article
            v-for="(t, i) in rcIntegrations"
            :key="i"
            class="rc-int-card rc-reveal"
            :style="{ transitionDelay: (i * 40) + 'ms' }"
          >
            <span v-if="t.tested" class="rc-int-tested"><span class="rc-int-d"></span> Testado</span>
            <div class="rc-int-row">
              <div class="rc-int-icon">{{ t.icon }}</div>
              <div class="rc-int-name">{{ t.name }}</div>
            </div>
            <div class="rc-int-where">{{ t.where }}</div>
            <div class="rc-int-env"><b>{{ t.env }}</b>=<span class="rc-int-url">router.cortexx.online</span></div>
          </article>
        </div>

        <div class="rc-int-more rc-reveal">
          <b>+ também roda em:</b> Zed · Vercel AI SDK · LlamaIndex · OpenRouter clients · Goose · Roo Code · Cody · scripts Python/Node — qualquer ferramenta que aceite uma chave compatível com OpenAI.
        </div>
      </section>

      <!-- Providers -->
      <section id="rc-providers" class="rc-section rc-wrap">
        <div class="rc-s-label"><span class="rc-sq"></span> $ providers <span class="rc-tic">--list --all</span></div>
        <h2 class="rc-s-title">Providers e modelos<br />liberados na chave.</h2>
        <p class="rc-s-sub">Tudo na chave certo? Tudo. Quando sai um modelo novo, a gente libera no mesmo dia — sem cobrar a mais.</p>
        <div class="rc-providers">
          <article v-for="pv in rcProviderCards" :key="pv.code" class="rc-prov-card">
            <span class="rc-prov-status"><span class="rc-prov-dot"></span> ON</span>
            <div class="rc-prov-ico" :style="{ color: pv.color, borderColor: pv.color + '33' }">{{ pv.code }}</div>
            <h4>{{ pv.name }}</h4>
            <div class="rc-prov-models">{{ pv.models }}</div>
          </article>
        </div>
      </section>

      <!-- Pricing -->
      <section id="rc-pricing" class="rc-section rc-wrap">
        <div class="rc-s-label"><span class="rc-sq"></span> $ pricing <span class="rc-tic">--currency=BRL</span></div>
        <h2 class="rc-s-title">Três planos.<br />Zero letra miúda.</h2>
        <p class="rc-s-sub">Cancela quando quiser. Pagamento facilitado via Pix ou cartão.</p>
        <div class="rc-pricing-grid">
          <article
            v-for="plan in rcPlans"
            :key="plan.name"
            class="rc-plan"
            :class="{ 'rc-plan-feat': plan.feat }"
          >
            <span v-if="plan.feat" class="rc-plan-badge">Recomendado</span>
            <div class="rc-plan-type">{{ plan.type }}</div>
            <h3>{{ plan.name }}</h3>
            <p class="rc-plan-desc">{{ plan.desc }}</p>
            <div class="rc-plan-price">
              <span class="rc-price-cur">R$</span>
              <span class="rc-price-amt">{{ plan.price }}</span>
              <span class="rc-price-per">/mês</span>
            </div>
            <div v-if="plan.oldPrice" class="rc-plan-old">Equivalente direto · R$ {{ plan.oldPrice }}/mês</div>
            <div v-else class="rc-plan-old rc-plan-old-spacer">·</div>
            <ul class="rc-plan-features">
              <li v-for="f in plan.features" :key="f">{{ f }}</li>
              <li v-for="f in plan.muted" :key="f" class="rc-feat-muted">{{ f }}</li>
            </ul>
            <router-link
              to="/purchase?tab=subscription"
              class="rc-btn rc-btn-lg"
              :class="plan.feat ? 'rc-btn-primary' : 'rc-btn-dark'"
              style="justify-content: center; display: flex;"
            >
              {{ plan.feat ? 'Assinar Pro' : 'Escolher ' + plan.name }}
              <span class="rc-arr">→</span>
            </router-link>
          </article>
        </div>
      </section>

      <!-- Proof -->
      <section id="rc-proof" class="rc-section rc-wrap">
        <div class="rc-s-label"><span class="rc-sq"></span> $ devs --aprovam <span class="rc-tic">--count=100+</span></div>
        <h2 class="rc-s-title">Devs brasileiros<br />respirando IA o dia inteiro.</h2>
        <p class="rc-s-sub">Quem usa Route-Cortexx tá fazendo refactor pesado, agentes autônomos e automações em loop. Não dá pra parar.</p>
        <div class="rc-proof-grid">
          <article v-for="(tm, i) in rcTestimonials" :key="i" class="rc-tcard">
            <div class="rc-tcard-meta">
              <span>{{ tm.meta[0] }}</span>
              <span class="rc-t-sep">/</span>
              <span class="rc-t-ec">{{ tm.meta[1] }}</span>
            </div>
            <p class="rc-tcard-quote">"{{ tm.quote }}"</p>
            <div class="rc-tcard-who">
              <div class="rc-tcard-av">{{ tm.initial }}</div>
              <div>
                <b>{{ tm.name }}</b>
                <span>{{ tm.role }}</span>
              </div>
            </div>
          </article>
        </div>
      </section>

      <!-- FAQ -->
      <section id="rc-faq" class="rc-section rc-wrap">
        <div class="rc-s-label"><span class="rc-sq"></span> $ faq <span class="rc-tic">--top=8</span></div>
        <h2 class="rc-s-title">Dúvidas que<br />aparecem no DM.</h2>
        <p class="rc-s-sub">Se a sua não tá aqui, chama no WhatsApp que respondo direto. Sem bot.</p>
        <div class="rc-faq-grid">
          <div
            v-for="(faq, i) in rcFaqs"
            :key="i"
            class="rc-faq-item"
            :class="{ 'rc-faq-open': openFaq === i }"
          >
            <button @click="toggleFaq(i)">
              <span>
                <span class="rc-faq-num">Q.{{ String(i + 1).padStart(2, '0') }}</span>
                {{ faq.q }}
              </span>
              <span class="rc-faq-plus">+</span>
            </button>
            <div class="rc-faq-answer"><p>{{ faq.a }}</p></div>
          </div>
        </div>
      </section>

      <!-- Final CTA -->
      <section class="rc-section rc-wrap">
        <div class="rc-final-cta">
          <div class="rc-s-label" style="justify-content: center;">
            <span class="rc-sq"></span> $ próximo_passo <span class="rc-tic">--ready</span>
          </div>
          <h2 class="rc-s-title">Pronto pra rotear seus tokens?</h2>
          <p class="rc-cta-p">Pix confirma em segundos. Em 4 minutos você tá rodando Claude Opus 4.7 no seu Cursor, pagando em real.</p>
          <div class="rc-cta-btns">
            <router-link to="/purchase?tab=subscription" class="rc-btn rc-btn-primary rc-btn-lg">
              Quero minha chave <span class="rc-arr">→</span>
            </router-link>
            <a :href="contactUrl" class="rc-btn rc-btn-ghost rc-btn-lg">Falar no WhatsApp</a>
          </div>
        </div>
      </section>

    </main>

    <!-- ── Footer ──────────────────────────────────────────── -->
    <footer class="rc-footer">
      <div class="rc-wrap rc-footer-grid">
        <div class="rc-footer-col rc-footer-about">
          <div class="rc-brand">
            <div class="rc-brand-logo">
              <img :src="siteLogo || '/logo.png'" alt="" />
            </div>
            <div class="rc-brand-name">
              <b>{{ siteName }}</b>
              <span>AI ROUTING · BR</span>
            </div>
          </div>
          <p class="rc-footer-desc">
            Gateway de roteamento e cota pra APIs de IA. Acesse Claude, GPT e Gemini
            com uma chave só, em real, no Pix. Não somos a Anthropic — somos seu intermediário brasileiro.
          </p>
        </div>
        <div class="rc-footer-col">
          <h5>Produto</h5>
          <a href="#rc-how">Como funciona</a>
          <a href="#rc-providers">Providers</a>
          <a href="#rc-pricing">Pricing</a>
          <router-link :to="dashboardPath">Painel</router-link>
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a>
        </div>
        <div class="rc-footer-col">
          <h5>Canais</h5>
          <a :href="contactUrl" target="_blank" rel="noopener noreferrer">WhatsApp</a>
          <a href="https://youtube.com" target="_blank" rel="noopener noreferrer">YouTube</a>
          <a href="https://instagram.com" target="_blank" rel="noopener noreferrer">Instagram</a>
        </div>
        <div class="rc-footer-col">
          <h5>Legal</h5>
          <a href="#">Termos de uso</a>
          <a href="#">Privacidade</a>
          <a href="#">Status</a>
          <a href="#">Reembolso</a>
        </div>
      </div>
      <div class="rc-wrap rc-footer-bot">
        <span>© {{ currentYear }} {{ siteName }} · Pagamento facilitado via Pix ou cartão</span>
        <span>route.cortexx.online</span>
      </div>
    </footer>

    <!-- ── Floating sticky CTA pill (aparece após scroll) ── -->
    <div v-if="!floatDismissed" class="rc-float-cta" :class="{ 'rc-float-visible': floatVisible }">
      <div class="rc-float-inner">
        <span class="rc-float-prompt">{{ '>_' }} Route-Cortexx</span>
        <span class="rc-float-sep">·</span>
        <a :href="contactUrl" target="_blank" rel="noopener noreferrer" class="rc-float-link">WhatsApp</a>
        <a href="#" class="rc-float-link">YouTube</a>
        <a href="#" class="rc-float-link">Instagram</a>
        <router-link to="/purchase?tab=subscription" class="rc-float-btn">Assinar · R$ 249 <span>→</span></router-link>
        <button class="rc-float-close" @click="floatDismissed = true" aria-label="Dispensar">×</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { useAuthStore, useAppStore } from '@/stores'
import { paymentAPI } from '@/api/payment'
import type { SubscriptionPlan } from '@/types/payment'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

// ── Auth & Settings ────────────────────────────────────────
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const currentYear = computed(() => new Date().getFullYear())

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Route-Cortexx')
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
const packetTick = ref(0)
interface LogEntry { id: number; t: string; path: string; prov: string; model: string; latency: number }
const logs = ref<LogEntry[]>([])
let routeInterval: ReturnType<typeof setInterval> | null = null

const rcRequests = [
  { y: 80, label: 'POST /v1/messages', method: 'POST /v1/messages', path: '/v1/messages', model: 'claude-opus-4.7', prov: 'AN_03' },
  { y: 140, label: 'POST /v1/chat/completions', method: 'POST /v1/chat/completions', path: '/v1/chat/completions', model: 'gpt-5', prov: 'OA_01' },
  { y: 200, label: 'POST /v1/embeddings', method: 'POST /v1/embeddings', path: '/v1/embeddings', model: 'gemini-2.5-flash', prov: 'GG_02' },
]
const rcProviders = [
  { y: 70, name: 'ANTHROPIC' },
  { y: 160, name: 'OPENAI' },
  { y: 250, name: 'GOOGLE' },
]

function tickRoute() {
  const next = (activeRoute.value + 1) % 3
  activeRoute.value = next
  packetTick.value++
  const now = new Date()
  const t = `${String(now.getHours()).padStart(2,'0')}:${String(now.getMinutes()).padStart(2,'0')}:${String(now.getSeconds()).padStart(2,'0')}`
  const latency = 280 + Math.floor(Math.random() * 220)
  const req = rcRequests[next]
  const entry: LogEntry = { id: now.getTime() + Math.random(), t, path: req.path, prov: req.prov, model: req.model, latency }
  logs.value = [entry, ...logs.value].slice(0, 3)
}

// Floating CTA pill state
const floatVisible = ref(false)
const floatDismissed = ref(false)
function onScroll() {
  if (floatDismissed.value) return
  floatVisible.value = window.scrollY > 700
}

// Tutorial video — COLE O LINK DO YOUTUBE AQUI quando estiver pronto
// ex.: 'https://www.youtube.com/embed/SEU_VIDEO_ID'
const videoEmbedUrl = ''

const rcChapters = [
  { t: '00:00', txt: 'Pagar via Pix em 3 cliques' },
  { t: '00:15', txt: 'Receber a chave no painel' },
  { t: '00:30', txt: 'Configurar no Cursor' },
  { t: '00:55', txt: 'Configurar no Claude Code' },
  { t: '01:15', txt: 'Primeira chamada · Opus 4.7' },
  { t: '01:25', txt: 'Monitorar uso e custo' },
]

const rcIntegrations = [
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
        e.target.classList.add('rc-reveal-in')
        revealObserver?.unobserve(e.target)
      }
    })
  }, { threshold: 0.12, rootMargin: '0px 0px -60px 0px' })
  document.querySelectorAll('.rc-reveal').forEach(el => revealObserver?.observe(el))
}

onMounted(() => {
  tickRoute()
  routeInterval = setInterval(tickRoute, 1400)
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
  loadPlans()
  // Defer reveal setup so DOM is fully painted
  setTimeout(setupReveal, 30)
})
onUnmounted(() => {
  if (routeInterval) clearInterval(routeInterval)
  window.removeEventListener('scroll', onScroll)
  if (revealObserver) revealObserver.disconnect()
})

// ── Math section ───────────────────────────────────────────
const activeMath = ref(0)
const rcMathRows = [
  { label: 'Claude Pro / Max', direct: 'R$ 580' },
  { label: 'ChatGPT Plus / Pro', direct: 'R$ 320' },
  { label: 'Gemini Advanced', direct: 'R$ 120' },
]
const rcMathPoints = [
  { tag: '01 · Fator', title: 'Compra em volume', body: 'Compramos pacotes de tokens dos providers em escala. Volume continuado virou desconto continuado — e a gente repassa pra você.' },
  { tag: '02 · Fator', title: 'Roteamento inteligente', body: 'Pra cada request, escolhemos o provider com a melhor latência e o menor custo no momento. Você paga sempre o melhor preço, sem decidir nada.' },
  { tag: '03 · Fator', title: 'Sem múltiplas assinaturas', body: 'Em vez de assinar Claude Max + ChatGPT Pro + Gemini Ultra (R$ 1.020+), você paga um plano só e usa todos. A mesma performance, em real.' },
]

// ── Providers ─────────────────────────────────────────────
const rcProviderCards = [
  { code: 'AN', name: 'Anthropic', models: 'Claude Opus 4.7 · Sonnet 4.6 · Haiku 4.5', color: '#d97757' },
  { code: 'OA', name: 'OpenAI', models: 'GPT-5 · o4 · gpt-realtime', color: '#10a37f' },
  { code: 'GG', name: 'Google', models: 'Gemini 2.5 Pro · Flash · Nano', color: '#4285f4' },
]

// ── Plans ─────────────────────────────────────────────────
interface DisplayPlan {
  type: string
  name: string
  feat: boolean
  price: number
  oldPrice: number | null
  desc: string
  features: string[]
  muted: string[]
}

const defaultPlans: DisplayPlan[] = [
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

const loadedPlans = ref<DisplayPlan[]>(defaultPlans)

function mapApiPlansToDisplay(apiPlans: SubscriptionPlan[]): DisplayPlan[] {
  const plans = apiPlans
    .filter(p => p.for_sale)
    .sort((a, b) => a.sort_order - b.sort_order)
    .map((p, i) => ({
      type: i === 1 ? 'Recomendado' : p.group_name || 'Plano',
      name: p.name,
      feat: i === 1,
      price: p.price,
      oldPrice: p.original_price || null,
      desc: p.description,
      features: typeof p.features === 'string' ? JSON.parse(p.features) : p.features,
      muted: [],
    }))
  return plans.length > 0 ? plans : defaultPlans
}

async function loadPlans() {
  try {
    const { data } = await paymentAPI.getPlans()
    if (data && data.length > 0) {
      loadedPlans.value = mapApiPlansToDisplay(data)
    }
  } catch {
    // Use default plans if API fails
  }
}

const rcPlans = computed(() => loadedPlans.value)

// ── Testimonials ──────────────────────────────────────────
const rcTestimonials = [
  { quote: 'Migrei 3 produtos pra Route-Cortexx num final de semana. Economia de R$ 2.300/mês e parou de cair quando a Anthropic engasga.', name: 'Marina V.', role: 'Tech lead · fintech SP', initial: 'MV', meta: ['ECONOMIA', 'R$ 2.300/mês'] },
  { quote: 'Setup foi o que demorou mais: 4 minutos. Troquei a base URL no Cursor e na hora estava rodando com Opus 4.7.', name: 'Caio P.', role: 'Indie hacker · POA', initial: 'CP', meta: ['SETUP', '4 MIN'] },
  { quote: 'O failover salvou meu deploy. Anthropic deu 529 quatro vezes esse mês e nenhum dos meus usuários percebeu.', name: 'Heitor M.', role: 'CTO · agência SaaS', initial: 'HM', meta: ['UPTIME', '99.94%'] },
]

// ── FAQ ───────────────────────────────────────────────────
const openFaq = ref(0)
function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? -1 : i
}
const rcFaqs = [
  { q: 'É a mesma API oficial dos providers?', a: 'Sim. O Route-Cortexx é só uma camada de roteamento e cota — as chamadas vão direto pros endpoints oficiais da Anthropic, OpenAI e Google. Sem proxy estranho. Você usa as SDKs oficiais, só troca a base URL e a chave.' },
  { q: 'Como é cobrado? Por token ou por mês?', a: 'Por mês, com pacote de tokens incluso. Você não precisa ficar somando dólar com câmbio nem se preocupar com pico de uso surpresa. Se passar do incluso, a gente avisa antes de cobrar mais.' },
  { q: 'Em qual ferramenta consigo usar?', a: 'Em qualquer coisa que aceite uma chave OpenAI-compat: Cursor, Claude Code, Cline, Aider, Continue, OpenRouter clients, n8n, scripts em Python/Node e por aí vai. Trocou a base URL e a chave, funcionou.' },
  { q: 'Tem garantia?', a: '7 dias. Se não rolar, devolvemos 100%, sem perguntinhas chatas. Pix ou no cartão.' },
  { q: 'Posso usar a mesma chave em vários PCs?', a: 'Pode. A chave é sua. Só não compartilhe fora do seu time porque a gente monitora padrão de uso e bloqueia abuso.' },
  { q: 'Vocês são oficiais da Anthropic/OpenAI?', a: 'Não, somos um intermediário brasileiro. Operamos sob os termos de uso de cada provider e somos parceiros revenda quando a categoria existe. Empresa registrada com sede no Brasil.' },
  { q: 'Quanto demora pra liberar o acesso?', a: 'Pix confirma em segundos. Cartão em até 2 minutos. A chave já fica no painel pronta pra colar.' },
  { q: 'E se um provider cair?', a: 'Failover automático. Configure no admin a ordem de prioridade e o Route-Cortexx redireciona pra próxima conta saudável sem você sentir.' },
]
</script>

<style scoped>
/* ── CSS custom properties ───────────────────────────────── */
.rc-home {
  --rc-bg:      #07090f;
  --rc-bg-1:    #0b0e16;
  --rc-bg-2:    #11151f;
  --rc-line:    rgba(255,255,255,0.08);
  --rc-line-2:  rgba(255,255,255,0.14);
  --rc-fg:      #e7eaf3;
  --rc-fg-1:    #b8bfd1;
  --rc-fg-2:    #8a93a8;
  --rc-fg-3:    #5b6378;
  --rc-cyan:    #5fd5ff;
  --rc-violet:  #a78bfa;
  --rc-emerald: #34d399;
  --rc-mono:    "JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, monospace;
  --rc-sans:    "Geist", ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif;
  --rc-r:       12px;

  background: var(--rc-bg);
  color: var(--rc-fg);
  font-family: var(--rc-sans);
  font-feature-settings: "ss01","cv11";
  -webkit-font-smoothing: antialiased;
  min-height: 100vh;
  overflow-x: hidden;
  position: relative;
}

/* ── Background layers ───────────────────────────────────── */
.rc-bg-grid {
  position: fixed; inset: 0; z-index: 0; pointer-events: none;
  background-image: radial-gradient(rgba(255,255,255,0.05) 1px, transparent 1px);
  background-size: 28px 28px;
  mask-image: radial-gradient(ellipse at 50% 30%, black 30%, transparent 80%);
}
.rc-bg-glow {
  position: fixed; inset: 0; z-index: 0; pointer-events: none;
  background:
    radial-gradient(700px 400px at 80% -10%, rgba(95,213,255,0.10), transparent 60%),
    radial-gradient(600px 400px at 0% 30%, rgba(167,139,250,0.08), transparent 60%);
}

/* ── Container ───────────────────────────────────────────── */
.rc-wrap { max-width: 1240px; margin: 0 auto; padding: 0 32px; position: relative; z-index: 1; }
@media (max-width: 720px) { .rc-wrap { padding: 0 20px; } }

/* ── Header ──────────────────────────────────────────────── */
.rc-header {
  position: sticky; top: 0; z-index: 50;
  backdrop-filter: blur(10px);
  background: rgba(7,9,15,0.72);
  border-bottom: 1px solid var(--rc-line);
}
.rc-header-inner { display: flex; align-items: center; justify-content: space-between; height: 64px; }
.rc-brand { display: flex; align-items: center; gap: 12px; text-decoration: none; }
.rc-brand-logo {
  width: 32px; height: 32px; border-radius: 8px; display: grid; place-items: center;
  background: #0e1320; border: 1px solid var(--rc-line-2); overflow: hidden;
}
.rc-brand-logo img { width: 28px; height: 28px; object-fit: contain; }
.rc-brand-name { display: flex; flex-direction: column; line-height: 1.05; }
.rc-brand-name b { font-weight: 600; font-size: 14px; letter-spacing: -0.01em; color: var(--rc-fg); }
.rc-brand-name span { font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3); letter-spacing: 0.04em; text-transform: uppercase; }
.rc-nav { display: flex; align-items: center; gap: 28px; }
@media (max-width: 860px) { .rc-nav { display: none; } }
.rc-nav-link { font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-2); letter-spacing: 0.04em; text-transform: uppercase; transition: color .15s; text-decoration: none; }
.rc-nav-link:hover { color: var(--rc-cyan); }
.rc-actions { display: flex; align-items: center; gap: 10px; }

/* ── Buttons ─────────────────────────────────────────────── */
.rc-btn {
  display: inline-flex; align-items: center; gap: 8px;
  height: 40px; padding: 0 16px; border-radius: 8px;
  font-family: var(--rc-mono); font-size: 12px; font-weight: 500;
  letter-spacing: 0.06em; text-transform: uppercase;
  border: 1px solid transparent; cursor: pointer; transition: all .15s;
  white-space: nowrap; text-decoration: none;
}
.rc-btn-lg { height: 48px; padding: 0 22px; font-size: 13px; }
.rc-btn-primary { background: var(--rc-cyan); color: #021018; border-color: var(--rc-cyan); }
.rc-btn-primary:hover { background: #82e1ff; border-color: #82e1ff; box-shadow: 0 0 0 4px rgba(95,213,255,0.15); }
.rc-btn-ghost { background: transparent; color: var(--rc-fg); border-color: var(--rc-line-2); }
.rc-btn-ghost:hover { border-color: var(--rc-cyan); color: var(--rc-cyan); }
.rc-btn-dark { background: #0e1320; color: var(--rc-fg); border-color: var(--rc-line-2); }
.rc-btn-dark:hover { background: #131a2a; }
.rc-arr { font-family: var(--rc-sans); font-weight: 500; }

/* ── Section primitives ──────────────────────────────────── */
.rc-section { padding: 96px 0; }
@media (max-width: 720px) { .rc-section { padding: 64px 0; } }
.rc-s-label {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--rc-mono); font-size: 11px; letter-spacing: 0.14em;
  color: var(--rc-cyan); text-transform: uppercase; margin-bottom: 16px;
}
.rc-sq { width: 6px; height: 6px; background: var(--rc-cyan); display: inline-block; }
.rc-tic { color: var(--rc-fg-3); }
.rc-s-title {
  font-size: clamp(34px, 4.2vw, 56px); font-weight: 600;
  line-height: 1.04; letter-spacing: -0.025em; margin: 0 0 18px;
  text-wrap: balance;
}
.rc-s-sub { font-size: 17px; color: var(--rc-fg-1); max-width: 640px; line-height: 1.55; margin: 0; }
.rc-accent { color: var(--rc-cyan); }
.rc-strike { text-decoration: line-through; color: var(--rc-fg-3); }

/* ── Hero ────────────────────────────────────────────────── */
.rc-hero { padding-top: 64px; padding-bottom: 80px; }
.rc-hero-grid { display: grid; grid-template-columns: 1.05fr 1fr; gap: 60px; align-items: center; }
@media (max-width: 1040px) { .rc-hero-grid { grid-template-columns: 1fr; gap: 48px; } }
.rc-eyebrow {
  display: inline-flex; align-items: center; gap: 10px;
  font-family: var(--rc-mono); font-size: 11px; text-transform: uppercase; letter-spacing: 0.14em;
  color: var(--rc-fg-1); background: #0e1320; border: 1px solid var(--rc-line-2);
  padding: 6px 12px; border-radius: 999px; margin-bottom: 18px;
}
.rc-eyebrow-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--rc-cyan); box-shadow: 0 0 12px var(--rc-cyan); }
.rc-eyebrow-v { color: var(--rc-violet); }
.rc-h1 { font-size: clamp(40px,5.8vw,76px); font-weight: 600; line-height: 0.98; letter-spacing: -0.03em; margin: 0 0 22px; text-wrap: balance; }
.rc-lede { font-size: 18px; line-height: 1.55; color: var(--rc-fg-1); max-width: 540px; margin: 0 0 32px; }
.rc-cta-row { display: flex; flex-wrap: wrap; gap: 12px; margin-bottom: 32px; }
.rc-pill-row { display: flex; flex-wrap: wrap; gap: 18px; font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-2); letter-spacing: 0.04em; }
.rc-pill-item { display: inline-flex; align-items: center; gap: 8px; }
.rc-dot-green { width: 8px; height: 8px; border-radius: 50%; background: var(--rc-emerald); box-shadow: 0 0 12px rgba(52,211,153,.6); display: inline-block; }
.rc-ok { color: var(--rc-emerald); }

/* ── Console ─────────────────────────────────────────────── */
.rc-console {
  background: linear-gradient(180deg, #0d1119 0%, #090c14 100%);
  border: 1px solid var(--rc-line-2); border-radius: var(--rc-r); overflow: hidden;
  box-shadow: 0 1px 0 rgba(255,255,255,.05) inset, 0 40px 80px -30px rgba(0,0,0,.6), 0 0 0 1px rgba(95,213,255,.06);
}
.rc-console-head {
  display: flex; align-items: center; gap: 10px; height: 36px; padding: 0 14px;
  border-bottom: 1px solid var(--rc-line); background: #0a0d15;
  font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3);
}
.rc-console-dots { display: inline-flex; gap: 6px; margin-right: 8px; }
.rc-console-dots i { width: 10px; height: 10px; border-radius: 50%; background: #1c2233; display: inline-block; }
.rc-console-title { color: var(--rc-fg-1); }
.rc-console-live { margin-left: auto; display: inline-flex; align-items: center; gap: 6px; color: var(--rc-emerald); }
.rc-live-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--rc-emerald); box-shadow: 0 0 10px var(--rc-emerald); display: inline-block; animation: rc-pulse 1.4s infinite; }
@keyframes rc-pulse { 0%,100% { opacity: 1; } 50% { opacity: .4; } }
.rc-console-body { padding: 22px; }

/* SVG diagram */
.rc-diagram { width: 100%; height: 300px; display: block; }
.rc-lane { fill: none; stroke: url(#rcG1); stroke-width: 1.5; opacity: .35; stroke-dasharray: 4 6; }
.rc-lane-active { opacity: 1; stroke-dasharray: 8 4; animation: rc-dash 1.2s linear infinite; }
@keyframes rc-dash { to { stroke-dashoffset: -36; } }
.rc-svg-label { font-family: var(--rc-mono); font-size: 10px; letter-spacing: .04em; }
.rc-svg-cyan { fill: #5fd5ff; }
.rc-svg-violet { fill: #a78bfa; }
.rc-svg-dim { fill: #5b6378; font-family: var(--rc-mono); font-size: 10px; }

.rc-console-meta { margin-top: 18px; display: grid; grid-template-columns: repeat(3,1fr); gap: 12px; font-family: var(--rc-mono); font-size: 11px; }
.rc-meta-cell { background: #0a0d15; border: 1px solid var(--rc-line); border-radius: 8px; padding: 10px 12px; }
.rc-meta-k { color: var(--rc-fg-3); letter-spacing: .08em; text-transform: uppercase; font-size: 10px; }
.rc-meta-v { color: var(--rc-fg); font-size: 16px; margin-top: 4px; letter-spacing: -.01em; font-feature-settings: "tnum"; }
.rc-meta-u { color: var(--rc-fg-3); font-size: 11px; margin-left: 2px; }

/* ── v2: Routing diagram packets & log feed ──────────────── */
.rc-packet { fill: var(--rc-cyan); filter: drop-shadow(0 0 6px var(--rc-cyan)); }
.rc-packet-1 { offset-path: path("M194 80 C 240 80, 240 160, 280 160 C 320 160, 320 70, 420 70"); animation: rc-packet-run 1.4s linear infinite; }
.rc-packet-2 { offset-path: path("M194 140 C 240 140, 240 160, 280 160 C 320 160, 320 160, 420 160"); animation: rc-packet-run 1.4s linear infinite; }
.rc-packet-3 { offset-path: path("M194 200 C 240 200, 240 160, 280 160 C 320 160, 320 250, 420 250"); animation: rc-packet-run 1.4s linear infinite; }
@keyframes rc-packet-run {
  0%   { offset-distance: 0%;  opacity: 0; }
  10%  { opacity: 1; }
  90%  { opacity: 1; }
  100% { offset-distance: 100%; opacity: 0; }
}
.rc-prov-svg { fill: #0a0d15; stroke: rgba(255,255,255,0.12); transition: stroke 0.25s, fill 0.25s; }
.rc-prov-svg-active { stroke: #5fd5ff; fill: #102232; }

.rc-routing-target {
  margin-top: 14px;
  display: flex; flex-wrap: wrap; align-items: center; gap: 8px;
  font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-2);
}
.rc-rt-label { color: var(--rc-fg-3); letter-spacing: 0.06em; text-transform: uppercase; }
.rc-rt-val {
  color: var(--rc-cyan);
  background: rgba(95, 213, 255, 0.08);
  border: 1px solid rgba(95, 213, 255, 0.2);
  padding: 3px 8px;
  border-radius: 4px;
  font-feature-settings: "tnum";
}

.rc-log-feed {
  margin-top: 14px;
  background: #06080d;
  border: 1px solid var(--rc-line);
  border-radius: 8px;
  padding: 10px 14px;
  font-family: var(--rc-mono);
  font-size: 11px;
  height: 84px;
  overflow: hidden;
  position: relative;
}
.rc-log-feed::before {
  content: "$ tail -f /var/log/cortexx/routes.log";
  display: block;
  font-size: 10px; color: var(--rc-fg-3);
  letter-spacing: 0.06em;
  padding-bottom: 6px;
  border-bottom: 1px dashed var(--rc-line);
  margin-bottom: 6px;
}
.rc-log-line {
  display: flex; gap: 8px; align-items: baseline;
  line-height: 1.5;
  animation: rc-log-in 0.4s ease;
  white-space: nowrap; overflow: hidden;
}
@keyframes rc-log-in {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: translateY(0); }
}
.rc-log-t { color: var(--rc-fg-3); }
.rc-log-method { color: var(--rc-violet); }
.rc-log-arr { color: var(--rc-fg-3); }
.rc-log-prov { color: var(--rc-cyan); }
.rc-log-model { color: var(--rc-fg-1); }
.rc-log-latency { color: var(--rc-fg-2); margin-left: auto; }
.rc-log-ok { color: var(--rc-emerald); }

/* ── Tutorial section ───────────────────────────────────── */
.rc-tut-grid { display: grid; grid-template-columns: 1.6fr 1fr; gap: 36px; margin-top: 48px; }
@media (max-width: 920px) { .rc-tut-grid { grid-template-columns: 1fr; gap: 28px; } }
.rc-video-wrap {
  position: relative;
  aspect-ratio: 16 / 9;
  width: 100%;
  border-radius: var(--rc-r);
  overflow: hidden;
  background: #0a0d15;
  border: 1px solid var(--rc-line-2);
  box-shadow: 0 1px 0 rgba(255,255,255,0.05) inset, 0 30px 60px -30px rgba(0,0,0,0.7);
}
.rc-video-wrap iframe { position: absolute; inset: 0; width: 100%; height: 100%; border: 0; }
.rc-video-placeholder {
  position: absolute; inset: 0;
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 18px;
  background:
    radial-gradient(600px 300px at 50% 60%, rgba(95,213,255,0.08), transparent 60%),
    repeating-linear-gradient(45deg, #0a0d15 0px, #0a0d15 12px, #0c1019 12px, #0c1019 24px);
  cursor: pointer;
  transition: background 0.2s;
}
.rc-video-placeholder:hover { background:
  radial-gradient(700px 350px at 50% 60%, rgba(95,213,255,0.14), transparent 60%),
  repeating-linear-gradient(45deg, #0a0d15 0px, #0a0d15 12px, #0c1019 12px, #0c1019 24px); }
.rc-video-badge {
  position: absolute; top: 18px; left: 18px;
  display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--rc-mono); font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase;
  color: var(--rc-cyan);
  background: rgba(7,9,15,0.85); border: 1px solid var(--rc-line-2);
  padding: 6px 12px; border-radius: 999px;
}
.rc-video-time {
  position: absolute; top: 18px; right: 18px;
  font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-2);
  background: rgba(7,9,15,0.85); border: 1px solid var(--rc-line-2);
  padding: 6px 12px; border-radius: 999px;
}
.rc-play-btn {
  width: 82px; height: 82px; border-radius: 50%;
  background: var(--rc-cyan); color: #021018;
  display: grid; place-items: center; border: none; cursor: pointer;
  animation: rc-play-pulse 2s ease infinite;
  transition: transform 0.15s;
}
.rc-play-btn:hover { transform: scale(1.06); }
.rc-play-btn svg { margin-left: 6px; }
@keyframes rc-play-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(95, 213, 255, 0.5); }
  50%      { box-shadow: 0 0 0 20px rgba(95, 213, 255, 0); }
}
.rc-video-hint { font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-2); letter-spacing: 0.08em; text-transform: uppercase; text-align: center; padding: 0 24px; }
.rc-video-hint code { color: var(--rc-cyan); }
.rc-video-strip {
  position: absolute; bottom: 16px; left: 18px; right: 18px;
  display: flex; gap: 14px; align-items: center; justify-content: space-between;
  font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3);
}
.rc-vs-left { display: inline-flex; align-items: center; gap: 8px; }
.rc-vs-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--rc-cyan); display: inline-block; }
.rc-vs-progress { flex: 1; max-width: 200px; height: 3px; background: rgba(255,255,255,0.1); border-radius: 2px; overflow: hidden; position: relative; }
.rc-vs-progress::after { content: ""; display: block; width: 22%; height: 100%; background: var(--rc-cyan); }

.rc-tut-h { font-family: var(--rc-mono); font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: var(--rc-fg-3); margin: 0 0 16px; }
.rc-tut-list { list-style: none; margin: 0; padding: 0; }
.rc-tut-list li {
  display: grid; grid-template-columns: 56px 1fr;
  align-items: baseline; gap: 14px;
  padding: 14px 0; border-top: 1px solid var(--rc-line);
  font-size: 15px; color: var(--rc-fg);
  transition: padding-left 0.15s, color 0.15s;
  cursor: pointer;
}
.rc-tut-list li:hover { padding-left: 8px; color: var(--rc-cyan); }
.rc-tut-list li:last-child { border-bottom: 1px solid var(--rc-line); }
.rc-tut-time { font-family: var(--rc-mono); font-size: 12px; color: var(--rc-cyan); letter-spacing: 0.04em; font-feature-settings: "tnum"; }

/* ── Integrations grid ─────────────────────────────────── */
.rc-int-grid { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; margin-top: 48px; }
@media (max-width: 1040px) { .rc-int-grid { grid-template-columns: repeat(2,1fr); } }
@media (max-width: 540px)  { .rc-int-grid { grid-template-columns: 1fr; } }
.rc-int-card {
  background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r);
  padding: 20px; position: relative;
  transition: border-color 0.2s, transform 0.2s, background 0.2s;
}
.rc-int-card:hover { border-color: var(--rc-cyan); background: linear-gradient(180deg, #0e1923 0%, #0a0d15 100%); transform: translateY(-2px); }
.rc-int-row { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.rc-int-icon { width: 36px; height: 36px; background: #06080d; border: 1px solid var(--rc-line-2); border-radius: 8px; display: grid; place-items: center; font-family: var(--rc-mono); font-size: 16px; font-weight: 600; color: var(--rc-cyan); }
.rc-int-name { font-size: 16px; font-weight: 600; letter-spacing: -0.01em; }
.rc-int-where { font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3); letter-spacing: 0.04em; text-transform: uppercase; margin-bottom: 10px; }
.rc-int-env { background: #06080d; border: 1px solid var(--rc-line); font-family: var(--rc-mono); font-size: 11px; padding: 8px 10px; border-radius: 6px; color: var(--rc-fg-2); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.rc-int-env b { color: var(--rc-cyan); font-weight: 500; }
.rc-int-url { color: var(--rc-violet); }
.rc-int-tested {
  position: absolute; top: 18px; right: 18px;
  font-family: var(--rc-mono); font-size: 9px; color: var(--rc-emerald);
  display: inline-flex; align-items: center; gap: 4px;
  letter-spacing: 0.1em; text-transform: uppercase;
}
.rc-int-d { width: 5px; height: 5px; border-radius: 50%; background: var(--rc-emerald); box-shadow: 0 0 6px var(--rc-emerald); display: inline-block; }
.rc-int-more {
  margin-top: 32px; padding: 20px 24px;
  border: 1px dashed var(--rc-line-2); border-radius: var(--rc-r);
  font-family: var(--rc-mono); font-size: 12px;
  color: var(--rc-fg-2); letter-spacing: 0.04em; line-height: 1.6;
}
.rc-int-more b { color: var(--rc-cyan); font-weight: 500; }

/* ── Floating CTA pill ───────────────────────────────────── */
.rc-float-cta {
  position: fixed; left: 50%; bottom: 20px;
  transform: translateX(-50%) translateY(120%);
  z-index: 80;
  transition: transform 0.35s cubic-bezier(0.5, 1.5, 0.3, 1);
  pointer-events: none;
}
.rc-float-visible { transform: translateX(-50%) translateY(0); pointer-events: auto; }
.rc-float-inner {
  display: flex; align-items: center; gap: 16px;
  padding: 8px 8px 8px 18px;
  background: rgba(11, 14, 22, 0.92);
  backdrop-filter: blur(16px);
  border: 1px solid var(--rc-line-2);
  border-radius: 999px;
  box-shadow: 0 24px 60px -20px rgba(0,0,0,0.7);
  font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-1);
}
.rc-float-prompt { color: var(--rc-cyan); font-weight: 600; letter-spacing: 0.02em; }
.rc-float-sep { color: var(--rc-line-2); margin: 0 -8px; }
.rc-float-link {
  color: var(--rc-fg-2); letter-spacing: 0.04em; text-transform: uppercase; font-size: 11px;
  transition: color 0.15s; text-decoration: none;
}
.rc-float-link:hover { color: var(--rc-cyan); }
.rc-float-btn {
  height: 36px; padding: 0 16px;
  background: var(--rc-cyan); color: #021018;
  border: none; border-radius: 999px;
  font-family: var(--rc-mono); font-size: 12px; font-weight: 500;
  letter-spacing: 0.04em; text-transform: uppercase;
  display: inline-flex; align-items: center; gap: 6px;
  cursor: pointer; text-decoration: none;
  transition: background 0.15s;
}
.rc-float-btn:hover { background: #82e1ff; }
.rc-float-close {
  width: 36px; height: 36px;
  display: grid; place-items: center;
  background: transparent; border: none;
  color: var(--rc-fg-3); cursor: pointer;
  border-radius: 50%;
  font-family: var(--rc-sans); font-size: 16px;
}
.rc-float-close:hover { color: var(--rc-fg); background: rgba(255,255,255,0.05); }
@media (max-width: 720px) {
  .rc-float-link, .rc-float-sep { display: none; }
}

/* ── Scroll reveal ───────────────────────────────────────── */
.rc-reveal {
  opacity: 0;
  transform: translateY(20px);
  transition: opacity 0.6s ease, transform 0.6s ease;
}
.rc-reveal-in { opacity: 1; transform: translateY(0); }
@media (prefers-reduced-motion: reduce) {
  .rc-reveal, .rc-reveal-in { opacity: 1; transform: none; transition: none; }
}

/* ── Steps ───────────────────────────────────────────────── */
.rc-steps { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; }
@media (max-width: 920px) { .rc-steps { grid-template-columns: 1fr; } }
.rc-step { background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r); padding: 28px; transition: border-color .2s; }
.rc-step:hover { border-color: var(--rc-line-2); }
.rc-step-num { font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-3); letter-spacing: .1em; }
.rc-step-num b { color: var(--rc-cyan); font-weight: 500; }
.rc-step h3 { font-size: 22px; font-weight: 600; letter-spacing: -.015em; margin: 14px 0 8px; }
.rc-step p { color: var(--rc-fg-1); font-size: 15px; line-height: 1.55; margin: 0 0 18px; }
.rc-codeblock { background: #06080d; border: 1px solid var(--rc-line); font-family: var(--rc-mono); font-size: 12px; padding: 12px 14px; border-radius: 8px; color: var(--rc-fg-1); overflow: hidden; }
.rc-c { color: var(--rc-fg-3); }
.rc-k { color: var(--rc-cyan); }
.rc-v { color: var(--rc-violet); }
.rc-s { color: var(--rc-emerald); }

/* ── Math ────────────────────────────────────────────────── */
.rc-math-grid { display: grid; grid-template-columns: 1.1fr 1fr; gap: 60px; align-items: center; margin-top: 56px; }
@media (max-width: 1000px) { .rc-math-grid { grid-template-columns: 1fr; gap: 40px; } }
.rc-compare { background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r); overflow: hidden; }
.rc-cmp-row { display: grid; grid-template-columns: 1.2fr 1fr 1fr; border-bottom: 1px solid var(--rc-line); font-family: var(--rc-mono); font-size: 13px; }
.rc-cmp-row:last-child { border-bottom: none; }
.rc-cmp-header { background: #0a0d15; }
.rc-cmp-header .rc-cmp-cell { font-size: 10px; letter-spacing: .1em; text-transform: uppercase; color: var(--rc-fg-3); padding: 14px 18px; }
.rc-cmp-header .rc-cmp-ours { color: var(--rc-cyan); }
.rc-cmp-cell { padding: 16px 18px; border-right: 1px solid var(--rc-line); color: var(--rc-fg-1); }
.rc-cmp-cell:last-child { border-right: none; }
.rc-cmp-label { color: var(--rc-fg-2); }
.rc-cmp-ours { color: var(--rc-fg); font-weight: 500; background: rgba(95,213,255,.03); }
.rc-cmp-total { background: #0a0d15; }
.rc-cmp-total .rc-cmp-cell { font-size: 15px; color: var(--rc-fg); padding: 18px; font-weight: 600; }
.rc-cmp-total .rc-cmp-ours { color: var(--rc-cyan); }
.rc-math-points { display: flex; flex-direction: column; gap: 24px; }
.rc-math-pt { border-left: 2px solid var(--rc-line-2); padding: 4px 0 4px 20px; cursor: default; }
.rc-math-pt-active { border-left-color: var(--rc-cyan); }
.rc-pt-tag { font-family: var(--rc-mono); font-size: 11px; letter-spacing: .1em; color: var(--rc-violet); text-transform: uppercase; }
.rc-math-pt h4 { font-size: 18px; font-weight: 600; letter-spacing: -.01em; margin: 6px 0 8px; }
.rc-math-pt p { color: var(--rc-fg-1); font-size: 14.5px; line-height: 1.55; margin: 0; }

/* ── Providers ───────────────────────────────────────────── */
.rc-providers { display: grid; grid-template-columns: repeat(3,1fr); gap: 16px; margin-top: 48px; }
@media (max-width: 860px) { .rc-providers { grid-template-columns: 1fr; } }
.rc-prov-card { background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r); padding: 22px; position: relative; transition: all .2s; }
.rc-prov-card:hover { border-color: var(--rc-line-2); transform: translateY(-2px); }
.rc-prov-status { position: absolute; top: 18px; right: 18px; font-family: var(--rc-mono); font-size: 10px; color: var(--rc-emerald); display: inline-flex; align-items: center; gap: 6px; }
.rc-prov-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--rc-emerald); box-shadow: 0 0 8px var(--rc-emerald); display: inline-block; }
.rc-prov-ico { width: 40px; height: 40px; border-radius: 10px; display: grid; place-items: center; font-family: var(--rc-mono); font-size: 16px; font-weight: 600; background: #0a0d15; border: 1px solid; margin-bottom: 14px; }
.rc-prov-card h4 { font-size: 16px; font-weight: 600; margin: 0 0 6px; letter-spacing: -.01em; }
.rc-prov-models { font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3); letter-spacing: .04em; line-height: 1.6; }

/* ── Pricing ─────────────────────────────────────────────── */
.rc-pricing-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; align-items: stretch; }
@media (max-width: 960px) { .rc-pricing-grid { grid-template-columns: 1fr; } }
.rc-plan { background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r); padding: 32px; position: relative; display: flex; flex-direction: column; }
.rc-plan-feat { border-color: var(--rc-cyan); background: linear-gradient(180deg, #0e1923 0%, #0a0d15 100%); }
.rc-plan-badge { position: absolute; top: -10px; left: 32px; font-family: var(--rc-mono); font-size: 10px; letter-spacing: .14em; text-transform: uppercase; background: var(--rc-cyan); color: #021018; padding: 4px 10px; border-radius: 4px; }
.rc-plan-type { font-family: var(--rc-mono); font-size: 11px; letter-spacing: .14em; text-transform: uppercase; color: var(--rc-fg-2); margin-bottom: 12px; }
.rc-plan h3 { font-size: 24px; font-weight: 600; letter-spacing: -.015em; margin: 0 0 8px; }
.rc-plan-desc { font-size: 14px; color: var(--rc-fg-2); line-height: 1.5; margin: 0 0 22px; min-height: 42px; }
.rc-plan-price { display: flex; align-items: baseline; gap: 4px; margin-bottom: 6px; }
.rc-price-cur { font-family: var(--rc-mono); font-size: 14px; color: var(--rc-fg-2); }
.rc-price-amt { font-size: 56px; font-weight: 600; letter-spacing: -.04em; line-height: 1; }
.rc-price-per { font-family: var(--rc-mono); font-size: 13px; color: var(--rc-fg-3); margin-left: 4px; }
.rc-plan-old { font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-3); text-decoration: line-through; margin-bottom: 22px; }
.rc-plan-old-spacer { text-decoration: none; visibility: hidden; }
.rc-plan-features { list-style: none; margin: 0 0 24px; padding: 0; flex: 1; }
.rc-plan-features li { font-size: 14px; color: var(--rc-fg-1); padding: 8px 0 8px 24px; position: relative; border-bottom: 1px solid var(--rc-line); }
.rc-plan-features li:last-child { border-bottom: none; }
.rc-plan-features li::before { content: "→"; position: absolute; left: 0; top: 8px; font-family: var(--rc-mono); color: var(--rc-cyan); font-size: 14px; }
.rc-plan-features li.rc-feat-muted { color: var(--rc-fg-3); }
.rc-plan-features li.rc-feat-muted::before { content: "·"; color: var(--rc-fg-3); }

/* ── Proof ───────────────────────────────────────────────── */
.rc-proof-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; }
@media (max-width: 1000px) { .rc-proof-grid { grid-template-columns: 1fr; } }
.rc-tcard { background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r); padding: 28px; display: flex; flex-direction: column; }
.rc-tcard-meta { font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3); margin-bottom: 16px; letter-spacing: .08em; text-transform: uppercase; display: flex; gap: 8px; }
.rc-t-sep { color: var(--rc-line-2); }
.rc-t-ec { color: var(--rc-emerald); }
.rc-tcard-quote { font-size: 17px; line-height: 1.5; letter-spacing: -.01em; color: var(--rc-fg); margin: 0 0 24px; text-wrap: pretty; flex: 1; }
.rc-tcard-who { display: flex; align-items: center; gap: 12px; }
.rc-tcard-av { width: 36px; height: 36px; border-radius: 50%; background: linear-gradient(135deg, var(--rc-cyan), var(--rc-violet)); display: grid; place-items: center; font-family: var(--rc-mono); font-size: 13px; font-weight: 600; color: #021018; flex-shrink: 0; }
.rc-tcard-who b { display: block; font-size: 14px; font-weight: 500; }
.rc-tcard-who span { display: block; font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3); }

/* ── FAQ ─────────────────────────────────────────────────── */
.rc-faq-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0 36px; margin-top: 32px; }
@media (max-width: 880px) { .rc-faq-grid { grid-template-columns: 1fr; } }
.rc-faq-item { border-top: 1px solid var(--rc-line); }
.rc-faq-item button {
  display: flex; align-items: center; justify-content: space-between; gap: 18px;
  width: 100%; padding: 20px 0; background: transparent; border: none; cursor: pointer; text-align: left;
  color: var(--rc-fg); font-size: 15.5px; font-weight: 500; letter-spacing: -.01em; font-family: inherit;
}
.rc-faq-item button:hover { color: var(--rc-cyan); }
.rc-faq-num { font-family: var(--rc-mono); font-size: 12px; color: var(--rc-fg-3); margin-right: 14px; }
.rc-faq-plus { color: var(--rc-cyan); font-family: var(--rc-mono); font-size: 18px; transition: transform .2s; flex-shrink: 0; }
.rc-faq-open .rc-faq-plus { transform: rotate(45deg); }
.rc-faq-answer { max-height: 0; overflow: hidden; transition: max-height .25s ease, padding .25s ease; }
.rc-faq-open .rc-faq-answer { max-height: 400px; padding: 0 0 24px 30px; }
.rc-faq-answer p { color: var(--rc-fg-1); font-size: 14.5px; line-height: 1.6; margin: 0; max-width: 56ch; }

/* ── Final CTA ───────────────────────────────────────────── */
.rc-final-cta {
  background: var(--rc-bg-1); border: 1px solid var(--rc-line); border-radius: var(--rc-r);
  padding: 56px; text-align: center; position: relative; overflow: hidden;
}
.rc-final-cta::before {
  content: ""; position: absolute; inset: 0; pointer-events: none;
  background: radial-gradient(600px 300px at 50% 0%, rgba(95,213,255,.12), transparent 70%);
}
.rc-final-cta .rc-s-label { display: inline-flex; }
.rc-final-cta .rc-s-title { font-size: clamp(34px,4vw,52px); position: relative; }
.rc-cta-p { font-size: 17px; color: var(--rc-fg-1); margin: 0 auto 28px; max-width: 520px; position: relative; }
.rc-cta-btns { display: inline-flex; gap: 12px; position: relative; flex-wrap: wrap; justify-content: center; }
@media (max-width: 720px) { .rc-final-cta { padding: 40px 24px; } }

/* ── Footer ──────────────────────────────────────────────── */
.rc-footer { border-top: 1px solid var(--rc-line); margin-top: 80px; padding: 48px 0 28px; position: relative; z-index: 1; }
.rc-footer-grid { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr; gap: 36px; }
@media (max-width: 880px) { .rc-footer-grid { grid-template-columns: 1fr 1fr; } }
.rc-footer-col h5 { font-family: var(--rc-mono); font-size: 11px; letter-spacing: .14em; text-transform: uppercase; color: var(--rc-fg-3); margin: 0 0 16px; }
.rc-footer-col a { display: block; font-size: 14px; color: var(--rc-fg-1); padding: 6px 0; text-decoration: none; }
.rc-footer-col a:hover { color: var(--rc-cyan); }
.rc-footer-desc { font-size: 14px; color: var(--rc-fg-2); line-height: 1.55; max-width: 360px; margin: 12px 0 0; }
.rc-footer-bot {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: 36px; padding-top: 24px; border-top: 1px solid var(--rc-line);
  font-family: var(--rc-mono); font-size: 11px; color: var(--rc-fg-3); letter-spacing: .04em;
}
@media (max-width: 720px) { .rc-footer-bot { flex-direction: column; gap: 12px; text-align: center; } }
</style>
