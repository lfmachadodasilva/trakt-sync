import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Input } from "./ui/input";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "./ui/form";
import { useForm, FormProvider } from "react-hook-form";
import { Button } from "./ui/button";

export function AccordionDemo() {
  const form = useForm();

  return (
    <Accordion
      type="single"
      className="w-full"
      defaultValue="item-1" // Ensures the first item is open by default
    >
      <AccordionItem value="item-1">
        <AccordionTrigger>trakt</AccordionTrigger>
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
                      <a
                        href="#"
                        className="btn btn-outline"
                        onClick={() => {
                          form.setValue("code", "new-secret-value");
                        }}
                      >
                        get code
                      </a>
                    </div>
                  </FormControl>
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
        <AccordionTrigger>emby</AccordionTrigger>
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
        <AccordionTrigger>Return Policy</AccordionTrigger>
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
  );
}
