"use strict";(self.webpackChunk=self.webpackChunk||[]).push([[842],{5977:(n,s,a)=>{a.r(s),a.d(s,{data:()=>p});const p={key:"v-b53ae1f0",path:"/guide/features/script.html",title:"Running scripts",lang:"en-US",frontmatter:{},excerpt:"",headers:[],filePathRelative:"guide/features/script.md",git:{updatedTime:1634291321e3,contributors:[{name:"Anton",email:"a.shumansky@gmail.com",commits:1}]}}},2544:(n,s,a)=>{a.r(s),a.d(s,{default:()=>e});const p=(0,a(6252).uE)('<h1 id="running-scripts" tabindex="-1"><a class="header-anchor" href="#running-scripts" aria-hidden="true">#</a> Running scripts</h1><p>This feature is used to run background scripts(cron jobs).</p><p>First you need to define script definition that implements hitrix.Script interface:</p><div class="language-go ext-go line-numbers-mode"><pre class="language-go"><code>\n<span class="token keyword">type</span> TestScript <span class="token keyword">struct</span> <span class="token punctuation">{</span><span class="token punctuation">}</span>\n\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Code</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">string</span> <span class="token punctuation">{</span>\n    <span class="token keyword">return</span> <span class="token string">&quot;test-script&quot;</span>\n<span class="token punctuation">}</span>\n\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Unique</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>\n    <span class="token comment">// if true you can&#39;t run more than one script at the same time</span>\n    <span class="token keyword">return</span> <span class="token boolean">false</span>\n<span class="token punctuation">}</span>\n\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Description</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">string</span> <span class="token punctuation">{</span>\n    <span class="token keyword">return</span> <span class="token string">&quot;script description&quot;</span>\n<span class="token punctuation">}</span>\n\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Run</span><span class="token punctuation">(</span>ctx context<span class="token punctuation">.</span>Context<span class="token punctuation">,</span> exit hitrix<span class="token punctuation">.</span>Exit<span class="token punctuation">)</span> <span class="token punctuation">{</span>\n    <span class="token comment">// put logic here</span>\n\t<span class="token keyword">if</span> shouldExitWithCode2 <span class="token punctuation">{</span>\n        exit<span class="token punctuation">.</span><span class="token function">Error</span><span class="token punctuation">(</span><span class="token punctuation">)</span>\t<span class="token comment">// you can exit script and specify exit code</span>\n    <span class="token punctuation">}</span>\n<span class="token punctuation">}</span>\n\n</code></pre><div class="line-numbers"><span class="line-number">1</span><br><span class="line-number">2</span><br><span class="line-number">3</span><br><span class="line-number">4</span><br><span class="line-number">5</span><br><span class="line-number">6</span><br><span class="line-number">7</span><br><span class="line-number">8</span><br><span class="line-number">9</span><br><span class="line-number">10</span><br><span class="line-number">11</span><br><span class="line-number">12</span><br><span class="line-number">13</span><br><span class="line-number">14</span><br><span class="line-number">15</span><br><span class="line-number">16</span><br><span class="line-number">17</span><br><span class="line-number">18</span><br><span class="line-number">19</span><br><span class="line-number">20</span><br><span class="line-number">21</span><br><span class="line-number">22</span><br><span class="line-number">23</span><br></div></div><p>Methods above are required. Optionally you can also implement these interfaces:</p><div class="language-go ext-go line-numbers-mode"><pre class="language-go"><code>\n<span class="token comment">// hitrix.ScriptInfinity interface</span>\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Infinity</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>\n    <span class="token comment">// run script and use blocking operation in cases you run all your code in goroutines</span>\n    <span class="token keyword">return</span> <span class="token boolean">true</span>\n<span class="token punctuation">}</span>\n\n<span class="token comment">// hitrix.ScriptInterval interface</span>\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Interval</span><span class="token punctuation">(</span><span class="token punctuation">)</span> time<span class="token punctuation">.</span>Duration <span class="token punctuation">{</span>                                                    \n    <span class="token comment">// run script every minute</span>\n    <span class="token keyword">return</span> time<span class="token punctuation">.</span>Minute \n<span class="token punctuation">}</span>\n\n<span class="token comment">// hitrix.ScriptIntervalOptional interface</span>\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">IntervalActive</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>                                                    \n    <span class="token comment">// only run first day of month</span>\n    <span class="token keyword">return</span> time<span class="token punctuation">.</span><span class="token function">Now</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Day</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token operator">==</span> <span class="token number">1</span>\n<span class="token punctuation">}</span>\n\n<span class="token comment">// hitrix.ScriptIntermediate interface</span>\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">IsIntermediate</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>                                                    \n    <span class="token comment">// script is intermediate, for example is listening for data in chain</span>\n    <span class="token keyword">return</span> <span class="token boolean">true</span>\n<span class="token punctuation">}</span>\n\n<span class="token comment">// hitrix.ScriptOptional interface</span>\n<span class="token keyword">func</span> <span class="token punctuation">(</span>script <span class="token operator">*</span>TestScript<span class="token punctuation">)</span> <span class="token function">Active</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>                                                    \n    <span class="token comment">// this script is visible only in local mode</span>\n    <span class="token keyword">return</span> <span class="token function">DIC</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">App</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">IsInLocalMode</span><span class="token punctuation">(</span><span class="token punctuation">)</span>\n<span class="token punctuation">}</span>\n\n</code></pre><div class="line-numbers"><span class="line-number">1</span><br><span class="line-number">2</span><br><span class="line-number">3</span><br><span class="line-number">4</span><br><span class="line-number">5</span><br><span class="line-number">6</span><br><span class="line-number">7</span><br><span class="line-number">8</span><br><span class="line-number">9</span><br><span class="line-number">10</span><br><span class="line-number">11</span><br><span class="line-number">12</span><br><span class="line-number">13</span><br><span class="line-number">14</span><br><span class="line-number">15</span><br><span class="line-number">16</span><br><span class="line-number">17</span><br><span class="line-number">18</span><br><span class="line-number">19</span><br><span class="line-number">20</span><br><span class="line-number">21</span><br><span class="line-number">22</span><br><span class="line-number">23</span><br><span class="line-number">24</span><br><span class="line-number">25</span><br><span class="line-number">26</span><br><span class="line-number">27</span><br><span class="line-number">28</span><br><span class="line-number">29</span><br><span class="line-number">30</span><br><span class="line-number">31</span><br></div></div><p>Once you defined script you can run it using RunScript method:</p><div class="language-go ext-go line-numbers-mode"><pre class="language-go"><code><span class="token keyword">package</span> main\n<span class="token keyword">import</span> <span class="token string">&quot;github.com/coretrix/hitrix&quot;</span>\n\n<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>\n\th <span class="token operator">:=</span> hitrix<span class="token punctuation">.</span><span class="token function">New</span><span class="token punctuation">(</span><span class="token string">&quot;app_name&quot;</span><span class="token punctuation">,</span> <span class="token string">&quot;your secret&quot;</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Build</span><span class="token punctuation">(</span><span class="token punctuation">)</span>\n\th<span class="token punctuation">.</span><span class="token function">RunBackgroundProcess</span><span class="token punctuation">(</span><span class="token keyword">func</span><span class="token punctuation">(</span>b <span class="token operator">*</span>hitrix<span class="token punctuation">.</span>BackgroundProcessor<span class="token punctuation">)</span> <span class="token punctuation">{</span>\n\t\tb<span class="token punctuation">.</span><span class="token function">RunScript</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>TestScript<span class="token punctuation">)</span>\n\t<span class="token punctuation">}</span><span class="token punctuation">)</span>\n<span class="token punctuation">}</span>\n</code></pre><div class="line-numbers"><span class="line-number">1</span><br><span class="line-number">2</span><br><span class="line-number">3</span><br><span class="line-number">4</span><br><span class="line-number">5</span><br><span class="line-number">6</span><br><span class="line-number">7</span><br><span class="line-number">8</span><br><span class="line-number">9</span><br></div></div><p>You can also register script as dynamic script and run it using program flag:</p><div class="language-go ext-go line-numbers-mode"><pre class="language-go"><code><span class="token keyword">package</span> main\n<span class="token keyword">import</span> <span class="token string">&quot;github.com/coretrix/hitrix&quot;</span>\n\n<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>\n\t\n    hitrix<span class="token punctuation">.</span><span class="token function">New</span><span class="token punctuation">(</span><span class="token string">&quot;app_name&quot;</span><span class="token punctuation">,</span> <span class="token string">&quot;your secret&quot;</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">RegisterDIService</span><span class="token punctuation">(</span>\n        <span class="token operator">&amp;</span>registry<span class="token punctuation">.</span>ServiceProvider<span class="token punctuation">{</span>\n            Name<span class="token punctuation">:</span>   <span class="token string">&quot;my-script&quot;</span><span class="token punctuation">,</span>\n            \n            Script<span class="token punctuation">:</span> <span class="token boolean">true</span><span class="token punctuation">,</span> <span class="token comment">// you need to set true here</span>\n            Build<span class="token punctuation">:</span> <span class="token keyword">func</span><span class="token punctuation">(</span>ctn di<span class="token punctuation">.</span>Container<span class="token punctuation">)</span> <span class="token punctuation">(</span><span class="token keyword">interface</span><span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token builtin">error</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>\n                <span class="token keyword">return</span> <span class="token operator">&amp;</span>TestScript<span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token boolean">nil</span>\n            <span class="token punctuation">}</span><span class="token punctuation">,</span>\n        <span class="token punctuation">}</span><span class="token punctuation">,</span>\n    <span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Build</span><span class="token punctuation">(</span><span class="token punctuation">)</span>\n<span class="token punctuation">}</span>\n</code></pre><div class="line-numbers"><span class="line-number">1</span><br><span class="line-number">2</span><br><span class="line-number">3</span><br><span class="line-number">4</span><br><span class="line-number">5</span><br><span class="line-number">6</span><br><span class="line-number">7</span><br><span class="line-number">8</span><br><span class="line-number">9</span><br><span class="line-number">10</span><br><span class="line-number">11</span><br><span class="line-number">12</span><br><span class="line-number">13</span><br><span class="line-number">14</span><br><span class="line-number">15</span><br><span class="line-number">16</span><br></div></div><p>You can see all available script by using special flag <strong>-list-scripts</strong>:</p><div class="language-bash ext-sh line-numbers-mode"><pre class="language-bash"><code>./app -list-scripts\n</code></pre><div class="line-numbers"><span class="line-number">1</span><br></div></div><p>To run script:</p><div class="language-bash ext-sh line-numbers-mode"><pre class="language-bash"><code>./app -run-script my-script\n</code></pre><div class="line-numbers"><span class="line-number">1</span><br></div></div>',14),t={},e=(0,a(3744).Z)(t,[["render",function(n,s){return p}]])},3744:(n,s)=>{s.Z=(n,s)=>{for(const[a,p]of s)n[a]=p;return n}}}]);