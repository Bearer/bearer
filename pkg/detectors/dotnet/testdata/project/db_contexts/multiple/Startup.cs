
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using FeelItaly.Models;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc.Authorization;
using System.Security.Claims;
using Microsoft.AspNetCore.Mvc.ViewFeatures;

namespace FeelItaly{

    public class Startup{

        public Startup(IConfiguration configuration){
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940
        public void ConfigureServices(IServiceCollection services){
            //var connection = @"Server=DESKTOP-GMKL1J1;Database=FeelItaly;Trusted_Connection=True;ConnectRetryCount=0";
            //var connection = @"Server=LAPTOP-TCFKJLM6;Database=FeelItaly;Trusted_Connection=True;ConnectRetryCount=0";
            //var connection = @"Server=DESKTOP-27V5D2G;Database=FeelItaly;Trusted_Connection=True;ConnectRetryCount=0";
            var connection = @"Server=localhost;Database=FeelItaly;User=sa;Password=mieimasters20LI4";
            services.AddDbContext<UtilizadorContext>(options => options.UseSqlServer(connection));

            services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
                    .AddCookie(options =>{
                        options.LoginPath = "/UtilizadorView/LoginUtilizador/";
                    });

            services.AddMvc();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IHostingEnvironment env){

            if (env.IsDevelopment()){
                app.UseDeveloperExceptionPage();
            }

            app.UseAuthentication();
            app.UseMvcWithDefaultRoute();
            app.UseStaticFiles();
        }

    }
}