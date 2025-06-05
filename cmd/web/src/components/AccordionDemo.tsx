import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Input } from "./ui/input";
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "./ui/form";
import { useForm, FormProvider } from "react-hook-form";
import { Button } from "./ui/button";
import { useState } from "react";

export function AccordionDemo() {
  const form = useForm();

  const [active, setActive] = useState("item-1");

  return (
    <>
      <Accordion type="single" className="w-full" value={active}>
        <AccordionItem value="item-1">
          <AccordionTrigger
            onClick={() => {
              setActive("item-1");
            }}
          >
            trakt
          </AccordionTrigger>
          <AccordionContent className="flex flex-col gap-4 text-balance w-[400px]">
            <FormProvider {...form}>
              <FormField
                control={form.control}
                name="clientId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>client id</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="client id"
                        autoComplete="new-password"
                        {...field}
                      />
                    </FormControl>
                    {/* <FormDescription>
                    This is your public display name.
                  </FormDescription> */}
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="clientSecret"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>client secret</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="client secret"
                        autoComplete="new-password"
                        {...field}
                      />
                    </FormControl>
                    {/* <FormDescription>
                    This is your public display name.
                  </FormDescription> */}
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="code"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>code</FormLabel>
                    <FormControl>
                      <div className="flex items-center gap-2">
                        <Input type="password" placeholder="code" {...field} />
                        <Button asChild className="btn btn-primary">
                          <a href="/login" target="_blank">
                            get code
                          </a>
                        </Button>
                      </div>
                    </FormControl>
                    <FormDescription className="text-sm text-left">
                      click the button to get the code. this will navigate you
                      to trakt and you will need to login to your account. after
                      that you will be shown a code that you need to copy and
                      paste here.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="redirectUrl"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>redirect url</FormLabel>
                    <FormControl>
                      <div className="flex items-center gap-2">
                        <Input
                          type="text"
                          placeholder="redirect url"
                          {...field}
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit">save</Button>
            </FormProvider>
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="item-2">
          <AccordionTrigger
            onClick={() => {
              setActive("item-2");
            }}
          >
            emby
          </AccordionTrigger>
          <AccordionContent className="flex flex-col gap-4 text-balance w-[400px]">
            <p>
              We offer worldwide shipping through trusted courier partners.
              Standard delivery takes 3-5 business days, while express shipping
              ensures delivery within 1-2 business days.
            </p>
            <p>
              All orders are carefully packaged and fully insured. Track your
              shipment in real-time through our dedicated tracking portal.
            </p>
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="item-3">
          <AccordionTrigger
            onClick={() => {
              setActive("item-3");
            }}
          >
            Return Policy
          </AccordionTrigger>
          <AccordionContent className="flex flex-col gap-4 text-balance w-[400px]">
            <p>
              We stand behind our products with a comprehensive 30-day return
              policy. If you&apos;re not completely satisfied, simply return the
              item in its original condition.
            </p>
            <p>
              Our hassle-free return process includes free return shipping and
              full refunds processed within 48 hours of receiving the returned
              item.
            </p>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
      <Button
        className="btn btn-primary mt-4"
        onClick={() => {
          setActive(active === "item-1" ? "item-2" : "item-1");
        }}
      >
        sync
      </Button>
    </>
  );
}
